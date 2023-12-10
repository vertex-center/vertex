package port

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	BaselinesAdapter interface {
		// GetLatest returns the latest available Baseline. This
		// will typically fetch the latest Baseline from a remote source.
		GetLatest(ctx context.Context, channel types.UpdatesChannel) (types.Baseline, error)
	}

	HardwareAdapter interface {
		Reboot(ctx context.Context) error
	}

	HardwareKernelAdapter interface {
		Reboot() error
	}

	SettingsAdapter interface {
		Get() (types.AdminSettings, error)
		SetChannel(channel types.UpdatesChannel) error
		SetWebhook(webhook string) error
	}

	SshAdapter interface {
		GetAll(ctx context.Context) ([]types.PublicKey, error)
		Add(ctx context.Context, key string, username string) error
		Remove(ctx context.Context, fingerprint string, username string) error
		GetUsers(ctx context.Context) ([]user.User, error)
	}

	SshKernelAdapter interface {
		GetAll(users []user.User) ([]types.PublicKey, error)
		Add(key string, user user.User) error
		Remove(fingerprint string, user user.User) error
		GetUsers() ([]user.User, error)
	}
)
