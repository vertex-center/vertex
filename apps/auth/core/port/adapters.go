package port

import "github.com/vertex-center/vertex/apps/auth/core/types"

type (
	AuthAdapter interface {
		CreateAccount(username string, credentials types.CredentialsArgon2id) error
		GetCredentials(login string) ([]types.CredentialsArgon2id, error)
		GetUsersByCredential(credentialID uint) ([]types.User, error)
		SaveSession(token *types.Session) error
		DeleteSession(token string) error
		GetSession(token string) (*types.Session, error)
		GetUser(username string) (types.User, error)
		GetUserByID(id uint) (types.User, error)
	}
)
