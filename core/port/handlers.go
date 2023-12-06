package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	AppsHandler interface {
		// Get handles the retrieval of all apps.
		Get(c *router.Context)
	}

	AuthHandler interface {
		// Login handles the login of a user.
		Login(c *router.Context)
		// Register handles the registration of a user.
		Register(c *router.Context)
		// Logout handles the logout of a user.
		Logout(c *router.Context)
	}

	DebugHandler interface {
		// HardReset do a hard reset of Vertex.
		HardReset(c *router.Context)
	}
)
