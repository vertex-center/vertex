package adapter

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (a *AuthDbAdapter) GetCredentials(login string) ([]types.CredentialsArgon2id, error) {
	var creds []types.CredentialsArgon2id
	err := a.db.Get().Find(&creds, &types.CredentialsArgon2id{Login: login}).Error
	return creds, err
}

func (a *AuthDbAdapter) SaveToken(token *types.Token) error {
	return a.db.Get().Transaction(func(tx *gorm.DB) error {
		err := tx.Omit(clause.Associations).Create(token).Error
		if err != nil {
			return err
		}
		*token = types.Token{Token: token.Token}
		return tx.Model(&types.Token{}).Preload("User").First(&token).Error
	})
}

func (a *AuthDbAdapter) RemoveToken(token string) error {
	return a.db.Get().Delete(&types.Token{Token: token}).Error
}

func (a *AuthDbAdapter) GetToken(token string) (*types.Token, error) {
	var t types.Token
	err := a.db.Get().Preload("User").First(&t, &types.Token{Token: token}).Error
	return &t, err
}
