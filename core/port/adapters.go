package port

import (
	"context"

	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	BaselinesAdapter interface {
		// GetLatest returns the latest available Baseline. This
		// will typically fetch the latest Baseline from a remote source.
		GetLatest(ctx context.Context, channel types.SettingsUpdatesChannel) (types.Baseline, error)
	}

	SettingsAdapter interface {
		GetSettings() types.Settings
		GetNotificationsWebhook() *string
		SetNotificationsWebhook(webhook string) error
		GetChannel() *types.SettingsUpdatesChannel
		SetChannel(channel types.SettingsUpdatesChannel) error
	}

	DataConfigAdapter interface {
		GetDataConfig() types.DataConfig
		GetDBMSName() types.DbmsName
		SetDBMSName(name types.DbmsName) error
	}

	SshAdapter interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string, username string) error
		Remove(fingerprint string, username string) error
		GetUsers() ([]user.User, error)
	}

	SshKernelAdapter interface {
		GetAll(users []user.User) ([]types.PublicKey, error)
		Add(key string, user user.User) error
		Remove(fingerprint string, user user.User) error
		GetUsers() ([]user.User, error)
	}
)
