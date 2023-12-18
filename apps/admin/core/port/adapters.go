package port

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/common/baseline"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	HardwareAdapter interface {
		Reboot(ctx context.Context) error
	}

	HardwareKernelAdapter interface {
		Reboot() error
	}

	SettingsAdapter interface {
		Get() (types.AdminSettings, error)
		SetChannel(channel baseline.Channel) error
		SetWebhook(webhook string) error
	}

	SshAdapter interface {
		GetAll(ctx context.Context) ([]types.PublicKey, error)
		Add(ctx context.Context, key string, username string) error
		Remove(ctx context.Context, fingerprint string, username string) error
		GetUsers(ctx context.Context) ([]user.User, error)
	}

	SshKernelAdapter interface {
		GetAll(ctx context.Context, users []user.User) ([]types.PublicKey, error)
		Add(ctx context.Context, key string, user user.User) error
		Remove(ctx context.Context, fingerprint string, user user.User) error
		GetUsers(ctx context.Context) ([]user.User, error)
	}
)
