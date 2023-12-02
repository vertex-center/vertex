package port

import (
	"context"

	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
)

type (
	AppsService interface {
		All() []app.Meta
	}

	DbService interface {
		GetCurrentDbms() types.DbmsName
		MigrateTo(dbms types.DbmsName) error
	}

	DebugService interface {
		HardReset()
	}

	AdminSettingsService interface {
		Get() (types.AdminSettings, error)
		Update(settings types.AdminSettings) error
		GetWebhook() (*string, error)
		SetWebhook(webhook string) error
		GetChannel() (types.UpdatesChannel, error)
		SetChannel(channel types.UpdatesChannel) error
	}

	ChecksService interface {
		CheckAll(ctx context.Context) <-chan types.CheckResponse
	}

	UpdateService interface {
		GetUpdate(channel types.UpdatesChannel) (*types.Update, error)
		InstallLatest(channel types.UpdatesChannel) error
	}
)
