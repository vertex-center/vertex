package port

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	MockAdminSettingsAdapter struct{ mock.Mock }
	MockBaselinesAdapter     struct{ mock.Mock }
	MockSshAdapter           struct{ mock.Mock }
	MockSshKernelAdapter     struct{ mock.Mock }
)

func (m *MockSshAdapter) GetAll(ctx context.Context) ([]types.PublicKey, error) {
	args := m.Called(ctx)
	return args.Get(0).([]types.PublicKey), args.Error(1)
}

func (m *MockSshAdapter) Add(ctx context.Context, key string, username string) error {
	args := m.Called(ctx, key, username)
	return args.Error(0)
}

func (m *MockSshAdapter) Remove(ctx context.Context, fingerprint string, username string) error {
	args := m.Called(ctx, fingerprint, username)
	return args.Error(0)
}

func (m *MockSshAdapter) GetUsers(ctx context.Context) ([]user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockSshKernelAdapter) GetAll(ctx context.Context, users []user.User) ([]types.PublicKey, error) {
	args := m.Called(ctx, users)
	return args.Get(0).([]types.PublicKey), args.Error(1)
}

func (m *MockSshKernelAdapter) Add(ctx context.Context, key string, user user.User) error {
	args := m.Called(ctx, key, user)
	return args.Error(0)
}

func (m *MockSshKernelAdapter) Remove(ctx context.Context, fingerprint string, user user.User) error {
	args := m.Called(ctx, fingerprint, user)
	return args.Error(0)
}

func (m *MockSshKernelAdapter) GetUsers(ctx context.Context) ([]user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.User), args.Error(1)
}
