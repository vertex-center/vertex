package port

import (
	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	MockSshAdapter struct {
		mock.Mock
	}

	MockSshKernelAdapter struct {
		mock.Mock
	}
)

func (m *MockSshAdapter) GetAll() ([]types.PublicKey, error) {
	args := m.Called()
	return args.Get(0).([]types.PublicKey), args.Error(1)
}

func (m *MockSshAdapter) Add(key string, username string) error {
	args := m.Called(key, username)
	return args.Error(0)
}

func (m *MockSshAdapter) Remove(fingerprint string, username string) error {
	args := m.Called(fingerprint, username)
	return args.Error(0)
}

func (m *MockSshAdapter) GetUsers() ([]user.User, error) {
	args := m.Called()
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockSshKernelAdapter) GetAll(users []user.User) ([]types.PublicKey, error) {
	args := m.Called(users)
	return args.Get(0).([]types.PublicKey), args.Error(1)
}

func (m *MockSshKernelAdapter) Add(key string, user user.User) error {
	args := m.Called(key, user)
	return args.Error(0)
}

func (m *MockSshKernelAdapter) Remove(fingerprint string, user user.User) error {
	args := m.Called(fingerprint, user)
	return args.Error(0)
}

func (m *MockSshKernelAdapter) GetUsers() ([]user.User, error) {
	args := m.Called()
	return args.Get(0).([]user.User), args.Error(1)
}
