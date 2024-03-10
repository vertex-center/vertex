package port

import (
	"github.com/vertex-center/vertex/server/apps/admin/core/types"
	"github.com/vertex-center/vertex/server/common/baseline"
)

type (
	SettingsAdapter interface {
		Get() (types.AdminSettings, error)
		SetChannel(channel baseline.Channel) error
		SetWebhook(webhook string) error
	}
)
