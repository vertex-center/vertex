package port

import (
	"context"

	"github.com/vertex-center/vertex/core/types"
)

type (
	MockBaselinesAdapter struct {
		GetLatestFunc  func(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error)
		GetLatestCalls int
	}

	MockSettingsAdapter struct {
		GetSettingsFunc              func() types.Settings
		GetSettingsCalls             int
		GetNotificationsWebhookFunc  func() *string
		GetNotificationsWebhookCalls int
		SetNotificationsWebhookFunc  func(webhook string) error
		SetNotificationsWebhookCalls int
		GetChannelFunc               func() *types.SettingsUpdatesChannel
		GetChannelCalls              int
		SetChannelFunc               func(channel types.SettingsUpdatesChannel) error
		SetChannelCalls              int
	}

	MockSshAdapter struct {
		GetAllFunc  func() ([]types.PublicKey, error)
		GetAllCalls int
		AddFunc     func(key string) error
		AddCalls    int
		RemoveFunc  func(fingerprint string) error
		RemoveCalls int
	}
)

func (m *MockBaselinesAdapter) GetLatest(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error) {
	m.GetLatestCalls++
	return m.GetLatestFunc(ctx, channel)
}

func (m *MockSettingsAdapter) GetSettings() types.Settings {
	m.GetSettingsCalls++
	return m.GetSettingsFunc()
}

func (m *MockSettingsAdapter) GetNotificationsWebhook() *string {
	m.GetNotificationsWebhookCalls++
	return m.GetNotificationsWebhookFunc()
}

func (m *MockSettingsAdapter) SetNotificationsWebhook(webhook string) error {
	m.SetNotificationsWebhookCalls++
	return m.SetNotificationsWebhookFunc(webhook)
}

func (m *MockSettingsAdapter) GetChannel() *types.SettingsUpdatesChannel {
	m.GetChannelCalls++
	return m.GetChannelFunc()
}

func (m *MockSettingsAdapter) SetChannel(channel types.SettingsUpdatesChannel) error {
	m.SetChannelCalls++
	return m.SetChannelFunc(channel)
}

func (m *MockSshAdapter) GetAll() ([]types.PublicKey, error) {
	m.GetAllCalls++
	return m.GetAllFunc()
}

func (m *MockSshAdapter) Add(key string) error {
	m.AddCalls++
	return m.AddFunc(key)
}

func (m *MockSshAdapter) Remove(fingerprint string) error {
	m.RemoveCalls++
	return m.RemoveFunc(fingerprint)
}
