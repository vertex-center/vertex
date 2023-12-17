package adapter

import (
	"time"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type emailDbAdapter struct {
	db storage.DB
}

func NewEmailDbAdapter(db storage.DB) *emailDbAdapter {
	return &emailDbAdapter{
		db: db,
	}
}

func (a *emailDbAdapter) CreateEmail(email *types.Email) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	// check that the email does not already exist for this user
	var count int
	err = tx.Get(&count, `
		SELECT COUNT(*)
		FROM emails
		WHERE user_id = $1 AND email = $2 AND deleted_at IS NULL
	`, email.UserID, email.Email)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if count > 0 {
		_ = tx.Rollback()
		return types.ErrEmailAlreadyExists
	}

	email.CreatedAt = time.Now().Unix()
	email.UpdatedAt = time.Now().Unix()

	_, err = tx.NamedExec(`
		INSERT INTO emails (id, user_id, email, created_at, updated_at)
		VALUES (:id, :user_id, :email, :created_at, :updated_at)
	`, *email)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (a *emailDbAdapter) GetEmails(userID uuid.UUID) ([]types.Email, error) {
	var emails []types.Email
	err := a.db.Select(&emails, `
		SELECT *
		FROM emails
		WHERE user_id = $1 AND deleted_at IS NULL
	`, userID)
	return emails, err
}

func (a *emailDbAdapter) DeleteEmail(userID uuid.UUID, email string) error {
	_, err := a.db.Exec(`
		UPDATE emails
		SET deleted_at = $1
		WHERE user_id = $2 AND email = $3
	`, time.Now().Unix(), userID, email)
	return err
}
