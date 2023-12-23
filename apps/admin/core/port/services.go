package port

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/common/baseline"
)

type (
	ChecksService interface {
		CheckAll(ctx context.Context) <-chan types.CheckResponse
	}

	DatabaseService interface {
		GetCurrentDbms() string
		MigrateTo(dbms string) error
	}

	SettingsService interface {
		Get() (types.AdminSettings, error)
		Update(settings types.AdminSettings) error
		GetWebhook() (*string, error)
		SetWebhook(webhook string) error
		GetChannel() (baseline.Channel, error)
		SetChannel(channel baseline.Channel) error
	}

	UpdateService interface {
		GetUpdate(channel baseline.Channel) (*types.Update, error)
	}
)
