package port

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	ChecksService interface {
		CheckAll(ctx context.Context) <-chan types.CheckResponse
	}

	DatabaseService interface {
		GetCurrentDbms() string
		MigrateTo(dbms string) error
	}

	HardwareService interface {
		GetHost() (types.Host, error)
		GetCPUs() ([]types.CPU, error)
		Reboot(ctx context.Context) error
	}

	HardwareKernelService interface {
		Reboot() error
	}

	SettingsService interface {
		Get() (types.AdminSettings, error)
		Update(settings types.AdminSettings) error
		GetWebhook() (*string, error)
		SetWebhook(webhook string) error
		GetChannel() (types.UpdatesChannel, error)
		SetChannel(channel types.UpdatesChannel) error
	}

	SshService interface {
		GetAll(ctx context.Context) ([]types.PublicKey, error)
		Add(ctx context.Context, key string, username string) error
		Delete(ctx context.Context, fingerprint string, username string) error
		GetUsers(ctx context.Context) ([]string, error)
	}

	SshKernelService interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string, username string) error
		Delete(fingerprint string, username string) error
		GetUsers() ([]user.User, error)
	}

	UpdateService interface {
		GetUpdate(channel types.UpdatesChannel) (*types.Update, error)
		InstallLatest(channel types.UpdatesChannel) error
	}
)
