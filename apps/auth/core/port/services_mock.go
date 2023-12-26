package port

import (
	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/apps/auth/core/types"
)

type MockAuthService struct{ mock.Mock }

var _ AuthService = (*MockAuthService)(nil)

func (m *MockAuthService) Login(login, password string) (types.Session, error) {
	args := m.Called(login, password)
	return args.Get(0).(types.Session), args.Error(1)
}

func (m *MockAuthService) Register(login, password string) (types.Session, error) {
	args := m.Called(login, password)
	return args.Get(0).(types.Session), args.Error(1)
}

func (m *MockAuthService) Logout(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthService) Verify(token string) (*types.Session, error) {
	args := m.Called(token)
	return nil, args.Error(0)
}
