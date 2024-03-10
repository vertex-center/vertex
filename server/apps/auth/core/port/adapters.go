package port

import (
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/auth/core/types"
)

type (
	AuthAdapter interface {
		CreateAccount(username string, credentials types.CredentialsArgon2id) error

		GetCredentials(login string) ([]types.CredentialsArgon2id, error)

		SaveSession(token *types.Session) error
		DeleteSession(token string) error
		GetSession(token string) (*types.Session, error)

		GetUser(username string) (types.User, error)
		GetUserByID(id uuid.UUID) (types.User, error)
		GetUsersByCredential(credentialID uuid.UUID) ([]types.User, error)
		PatchUser(user types.User) (types.User, error)
		GetUserCredentialsMethods(userID uuid.UUID) ([]types.CredentialsMethods, error)
	}

	EmailAdapter interface {
		CreateEmail(email *types.Email) error
		GetEmails(userID uuid.UUID) ([]types.Email, error)
		DeleteEmail(userID uuid.UUID, email string) error
	}
)
