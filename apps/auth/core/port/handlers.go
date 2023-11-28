package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	AuthHandler interface {
		// Login handles the login of a user.
		Login(c *router.Context)
		// Register handles the registration of a user.
		Register(c *router.Context)
		// Logout handles the logout of a user.
		Logout(c *router.Context)
	}
)
