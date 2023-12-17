package adapter

import (
	"fmt"
	"time"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

type authDbAdapter struct {
	db storage.DB
}

func NewAuthDbAdapter(db storage.DB) port.AuthAdapter {
	return &authDbAdapter{db}
}

func (a *authDbAdapter) CreateAccount(username string, credential types.CredentialsArgon2id) error {
	user := types.User{
		ID:        uuid.New(),
		Username:  username,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		INSERT INTO users (id, username, created_at, updated_at)
		VALUES (:id, :username, :created_at, :updated_at)
	`, user)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	credential.ID = uuid.New()
	credential.CreatedAt = time.Now().Unix()
	credential.UpdatedAt = time.Now().Unix()

	_, err = tx.NamedExec(`
		INSERT INTO credentials_argon2 (id, login, hash, type, iterations, memory, parallelism, salt, key_len, created_at, updated_at)
		VALUES (:id, :login, :hash, :type, :iterations, :memory, :parallelism, :salt, :key_len, :created_at, :updated_at)
	`, credential)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO credentials_argon2_users (credential_id, user_id)
		VALUES ($1, $2)
	`, credential.ID, user.ID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	return tx.Commit()
}

func (a *authDbAdapter) GetCredentials(login string) ([]types.CredentialsArgon2id, error) {
	var creds []types.CredentialsArgon2id

	rows, err := a.db.Queryx(`
		SELECT *
		FROM credentials_argon2
		WHERE credentials_argon2.login = $1
	`, login)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cred types.CredentialsArgon2id
		err := rows.StructScan(&cred)
		if err != nil {
			return nil, err
		}
		creds = append(creds, cred)
	}

	return creds, nil
}

func (a *authDbAdapter) GetUsersByCredential(credentialID uuid.UUID) ([]types.User, error) {
	var users []types.User
	err := a.db.Select(&users, `
		SELECT users.id, users.username, users.created_at, users.updated_at, users.deleted_at
		FROM users
		JOIN credentials_argon2_users ON users.id = credentials_argon2_users.user_id
		WHERE credentials_argon2_users.credential_id = $1
	`, credentialID)
	return users, err
}

func (a *authDbAdapter) SaveSession(session *types.Session) error {
	session.CreatedAt = time.Now().Unix()
	session.UpdatedAt = time.Now().Unix()

	_, err := a.db.NamedExec(`
		INSERT INTO sessions (id, token, user_id, created_at, updated_at)
		VALUES (:id, :token, :user_id, :created_at, :updated_at)
	`, session)
	return err
}

func (a *authDbAdapter) DeleteSession(token string) error {
	_, err := a.db.Exec("UPDATE sessions SET deleted_at = $1 WHERE token = $2", time.Now().Unix(), token)
	return err
}

func (a *authDbAdapter) GetSession(token string) (*types.Session, error) {
	var session types.Session
	err := a.db.Get(&session, `
		SELECT *
		FROM sessions
		WHERE sessions.deleted_at IS NULL AND sessions.token = $1
	`, token)
	return &session, err
}

func (a *authDbAdapter) GetUser(username string) (types.User, error) {
	var user types.User
	err := a.db.Get(&user, `
		SELECT *
		FROM users
		WHERE users.username = $1
	`, username)
	return user, err
}

func (a *authDbAdapter) GetUserByID(id uuid.UUID) (types.User, error) {
	var user types.User
	err := a.db.Get(&user, `
		SELECT *
		FROM users
		WHERE users.id = $1
	`, id)
	return user, err
}

func (a *authDbAdapter) PatchUser(user types.User) (types.User, error) {
	user.UpdatedAt = time.Now().Unix()

	tx, err := a.db.Beginx()
	if err != nil {
		return user, err
	}

	_, err = tx.NamedExec(`
		UPDATE users
		SET username = :username, updated_at = :updated_at
		WHERE id = :id
	`, user)
	if err != nil {
		_ = tx.Rollback()
		return user, err
	}

	err = tx.Get(&user, `
		SELECT *
		FROM users
		WHERE users.id = $1
	`, user.ID)
	if err != nil {
		_ = tx.Rollback()
		return user, err
	}

	return user, tx.Commit()
}

func (a *authDbAdapter) GetUserCredentialsMethods(userID uuid.UUID) ([]types.CredentialsMethods, error) {
	var credsArgon2id []types.CredentialsArgon2id
	err := a.db.Select(&credsArgon2id, `
		SELECT credentials_argon2.*
		FROM credentials_argon2
		JOIN credentials_argon2_users ON credentials_argon2.id = credentials_argon2_users.credential_id
		WHERE credentials_argon2_users.user_id = $1
	`, userID)

	var creds []types.CredentialsMethods
	for _, cr := range credsArgon2id {
		creds = append(creds, types.CredentialsMethods{
			Name:        types.CredentialsTypeLoginPassword,
			Description: fmt.Sprintf("Your password is hashed with the %s algorithm", cr.Type),
		})
	}

	return creds, err
}
