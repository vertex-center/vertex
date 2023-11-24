package port

import (
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	AppsService interface {
		All() []app.Meta
	}

	DebugService interface {
		HardReset()
	}

	HardwareService interface {
		Get() types.Hardware
	}

	SettingsService interface {
		Get() types.Settings
		Update(settings types.Settings) error
		GetNotificationsWebhook() *string
		SetNotificationsWebhook(webhook string) error
		GetChannel() types.SettingsUpdatesChannel
		SetChannel(channel types.SettingsUpdatesChannel) error
	}

	SshService interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string) error
		Delete(fingerprint string) error
		GetUsers() ([]user.User, error)
	}

	UpdateService interface {
		GetUpdate(channel types.SettingsUpdatesChannel) (*types.Update, error)
		InstallLatest(channel types.SettingsUpdatesChannel) error
	}
)
