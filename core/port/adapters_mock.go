package port

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	MockBaselinesAdapter struct {
		mock.Mock
	}

	MockSettingsAdapter struct {
		mock.Mock
	}

	MockSshAdapter struct {
		mock.Mock
	}
)

func (m *MockBaselinesAdapter) GetLatest(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error) {
	args := m.Called(ctx, channel)
	return args.Get(0).(types.Baseline), args.Error(1)
}

func (m *MockSettingsAdapter) GetSettings() types.Settings {
	args := m.Called()
	return args.Get(0).(types.Settings)
}

func (m *MockSettingsAdapter) GetNotificationsWebhook() *string {
	args := m.Called()
	return args.Get(0).(*string)
}

func (m *MockSettingsAdapter) SetNotificationsWebhook(webhook string) error {
	args := m.Called(webhook)
	return args.Error(0)
}

func (m *MockSettingsAdapter) GetChannel() *types.SettingsUpdatesChannel {
	args := m.Called()
	return args.Get(0).(*types.SettingsUpdatesChannel)
}

func (m *MockSettingsAdapter) SetChannel(channel types.SettingsUpdatesChannel) error {
	args := m.Called(channel)
	return args.Error(0)
}

func (m *MockSshAdapter) GetAll() ([]types.PublicKey, error) {
	args := m.Called()
	return args.Get(0).([]types.PublicKey), args.Error(1)
}

func (m *MockSshAdapter) Add(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockSshAdapter) Remove(fingerprint string) error {
	args := m.Called(fingerprint)
	return args.Error(0)
}

func (m *MockSshAdapter) GetUsers() ([]user.User, error) {
	args := m.Called()
	return args.Get(0).([]user.User), args.Error(1)
}
