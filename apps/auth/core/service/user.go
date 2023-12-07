package service

import (
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
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

func (s *userService) GetUserByID(id uint) (types.User, error) {
	return s.adapter.GetUserByID(id)
}

func (s *userService) PatchUser(user types.User) (types.User, error) {
	return s.adapter.PatchUser(user)
}

func (s *userService) GetUserCredentialsMethods(userID uint) ([]types.CredentialsMethods, error) {
	return s.adapter.GetUserCredentialsMethods(userID)
}
