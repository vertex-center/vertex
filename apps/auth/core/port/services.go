package port

import (
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/auth/core/types"
)

type (
	AuthService interface {
		Login(login, password string) (types.Session, error)
		Register(login, password string) (types.Session, error)
		Logout(token string) error
		Verify(token string) (*types.Session, error)
	}

	EmailService interface {
		CreateEmail(userID uuid.UUID, email string) (types.Email, error)
		GetEmails(userID uuid.UUID) ([]types.Email, error)
		DeleteEmail(userID uuid.UUID, email string) error
	}

	MigrationService interface{}

	UserService interface {
		GetUser(username string) (types.User, error)
		GetUserByID(id uuid.UUID) (types.User, error)
		PatchUser(user types.User) (types.User, error)
		GetUserCredentialsMethods(userID uuid.UUID) ([]types.CredentialsMethods, error)
	}
)
