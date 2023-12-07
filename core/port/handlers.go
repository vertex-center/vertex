package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	AppsHandler interface {
		// Get handles the retrieval of all apps.
		Get(c *router.Context)
	}

	DebugHandler interface {
		// HardReset do a hard reset of Vertex.
		HardReset(c *router.Context)
	}
)
