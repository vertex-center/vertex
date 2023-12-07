package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	AppsHandler interface {
		GetApps(c *router.Context)
		GetAppsInfo() []oapi.Info
	}

	DebugHandler interface {
		HardReset(c *router.Context)
		HardResetInfo() []oapi.Info
	}
)
