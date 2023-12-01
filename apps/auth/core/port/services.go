package port

import "github.com/vertex-center/vertex/apps/auth/core/types"

type (
	AuthService interface {
		Login(login, password string) (types.Token, error)
		Register(login, password string) (types.Token, error)
		Logout(token string) error
		Verify(token string) error
	}

	MigrationService interface{}
)