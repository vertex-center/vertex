package port

import (
	"context"

	"github.com/vertex-center/vertex/core/types"
)

type (
	BaselinesAdapter interface {
		// GetLatest returns the latest available Baseline. This
		// will typically fetch the latest Baseline from a remote source.
		GetLatest(ctx context.Context, channel types.UpdatesChannel) (types.Baseline, error)
	}

	AdminSettingsAdapter interface {
		Get() (types.AdminSettings, error)
		SetChannel(channel types.UpdatesChannel) error
		SetWebhook(webhook string) error
	}

	DbAdapter interface {
		Get() *types.DB
		Connect() error
		GetDbConfig() types.DbConfig
		GetDBMSName() types.DbmsName
		SetDBMSName(name types.DbmsName) error
	}
)
