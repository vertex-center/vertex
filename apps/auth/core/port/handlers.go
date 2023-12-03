package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	AuthHandler interface {
		Login(c *router.Context)
		Register(c *router.Context)
		Logout(c *router.Context)
	}

	EmailHandler interface {
		CreateCurrentUserEmail(c *router.Context)
		GetCurrentUserEmails(c *router.Context)
		DeleteCurrentUserEmail(c *router.Context)
	}

	UserHandler interface {
		GetCurrentUser(c *router.Context)
		PatchCurrentUser(c *router.Context)
		GetCurrentUserCredentials(c *router.Context)
	}
)
