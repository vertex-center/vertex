package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	AuthHandler interface {
		Login(c *router.Context)
		Register(c *router.Context)
		Logout(c *router.Context)
	}

	UserHandler interface {
		GetCurrentUser(c *router.Context)
		PatchCurrentUser(c *router.Context)
		GetCurrentUserCredentials(c *router.Context)
	}
)
