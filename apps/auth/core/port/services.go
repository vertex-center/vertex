package port

import "github.com/vertex-center/vertex/apps/auth/core/types"

type (
	AuthService interface {
		Login(login, password string) (types.Session, error)
		Register(login, password string) (types.Session, error)
		Logout(token string) error
		Verify(token string) (*types.Session, error)
	}

	MigrationService interface{}

	UserService interface {
		GetUser(username string) (types.User, error)
		GetUserByID(id uint) (types.User, error)
		PatchUser(user types.User) (types.User, error)
		GetUserCredentialsMethods(userID uint) ([]types.CredentialsMethods, error)
	}
)
