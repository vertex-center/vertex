package port

import "github.com/vertex-center/vertex/apps/auth/core/types"

type (
	AuthAdapter interface {
		CreateAccount(username string, credentials types.CredentialsArgon2id) error
		GetCredentials(login string) ([]types.CredentialsArgon2id, error)
		SaveToken(token *types.Token) error
		RemoveToken(token string) error
		GetToken(token string) (*types.Token, error)
	}
)
