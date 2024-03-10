package service

import (
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/auth/core/port"
	"github.com/vertex-center/vertex/server/apps/auth/core/types"
)

type userService struct {
	adapter port.AuthAdapter
}

func NewUserService(adapter port.AuthAdapter) port.UserService {
	return &userService{
		adapter: adapter,
	}
}

func (s *userService) GetUser(username string) (types.User, error) {
	return s.adapter.GetUser(username)
}

func (s *userService) GetUserByID(id uuid.UUID) (types.User, error) {
	return s.adapter.GetUserByID(id)
}

func (s *userService) PatchUser(user types.User) (types.User, error) {
	return s.adapter.PatchUser(user)
}

func (s *userService) GetUserCredentialsMethods(userID uuid.UUID) ([]types.CredentialsMethods, error) {
	return s.adapter.GetUserCredentialsMethods(userID)
}
