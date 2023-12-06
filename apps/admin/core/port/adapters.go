package port

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	AdminSettingsAdapter interface {
		Get() (types.AdminSettings, error)
		SetChannel(channel types.UpdatesChannel) error
		SetWebhook(webhook string) error
	}

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
