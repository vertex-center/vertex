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

func (m *MockAuthAdapter) GetUsersByCredential(credentialID uint) ([]types.User, error) {
	args := m.Called(credentialID)
	return args.Get(0).([]types.User), args.Error(1)
}

func (m *MockAuthAdapter) SaveSession(session *types.Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockAuthAdapter) DeleteSession(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthAdapter) GetSession(token string) (*types.Session, error) {
	args := m.Called(token)
	return args.Get(0).(*types.Session), args.Error(1)
}

func (m *MockAuthAdapter) GetUser(username string) (types.User, error) {
	args := m.Called(username)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockAuthAdapter) GetUserByID(id uint) (types.User, error) {
	args := m.Called(id)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockAuthAdapter) PatchUser(user types.User) (types.User, error) {
	args := m.Called(user)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockAuthAdapter) GetUserCredentialsMethods(userID uint) ([]types.CredentialsMethods, error) {
	args := m.Called(userID)
	return args.Get(0).([]types.CredentialsMethods), args.Error(1)
}
