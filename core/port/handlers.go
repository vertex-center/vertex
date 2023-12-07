package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	AppsHandler interface {
		GetApps(c *router.Context)
	}

	DebugHandler interface {
		HardReset(c *router.Context)
	}
)
