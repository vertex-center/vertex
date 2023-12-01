package port

import (
	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/apps/auth/core/types"
)

type MockAuthAdapter struct {
	mock.Mock
}

func (m *MockAuthAdapter) CreateAccount(username string, credentials types.CredentialsArgon2id) error {
	args := m.Called(username, credentials)
	return args.Error(0)
}

func (m *MockAuthAdapter) GetCredentials(login string) ([]types.CredentialsArgon2id, error) {
	args := m.Called(login)
	return args.Get(0).([]types.CredentialsArgon2id), args.Error(1)
}

func (m *MockAuthAdapter) SaveToken(token *types.Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthAdapter) RemoveToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthAdapter) GetToken(token string) (*types.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*types.Token), args.Error(1)
}

func (m *MockAuthAdapter) GetUser(username string) (types.User, error) {
	args := m.Called(username)
	return args.Get(0).(types.User), args.Error(1)
}
