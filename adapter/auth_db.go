package adapter

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"gorm.io/gorm"
)

type AuthDbAdapter struct {
	db port.DbConfigAdapter
}

func NewAuthDbAdapter(db port.DbConfigAdapter) port.AuthAdapter {
	return &AuthDbAdapter{
		db: db,
	}
}

func (a *AuthDbAdapter) CreateAccount(username string, credentials types.CredentialsArgon2id) error {
	user := types.User{
		Username: username,
	}

	return a.db.Get().Transaction(func(tx *gorm.DB) error {
		// Caution: never change this to FirstOrCreate, as it will allow
		// anyone to take over an existing account by simply registering
		// with the same username.
		// To add the ability to create new credentials for an existing
		// user, we need to create a separate endpoint.
		err := tx.Create(&user).Error
		if err != nil {
			return err
		}
		return tx.Model(&user).Association("CredentialsArgon2id").Append(&credentials)
	})
}
