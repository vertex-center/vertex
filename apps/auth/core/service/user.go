package service

import (
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
)

type UserService struct {
	adapter port.AuthAdapter
}

func NewUserService(adapter port.AuthAdapter) port.UserService {
	return &UserService{
		adapter: adapter,
	}
}

func (s *UserService) GetUser(username string) (types.User, error) {
	return s.adapter.GetUser(username)
}

func (s *UserService) GetUserByID(id uint) (types.User, error) {
	return s.adapter.GetUserByID(id)
}

func (s *UserService) PatchUser(user types.User) (types.User, error) {
	return s.adapter.PatchUser(user)
}
