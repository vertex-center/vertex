package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	AuthHandler interface {
		Login(c *router.Context)
		LoginInfo() []oapi.Info

		Register(c *router.Context)
		RegisterInfo() []oapi.Info

		Logout(c *router.Context)
		LogoutInfo() []oapi.Info

		Verify(c *router.Context)
		VerifyInfo() []oapi.Info
	}

	EmailHandler interface {
		CreateCurrentUserEmail(c *router.Context)
		CreateCurrentUserEmailInfo() []oapi.Info

		GetCurrentUserEmails(c *router.Context)
		GetCurrentUserEmailsInfo() []oapi.Info

		DeleteCurrentUserEmail(c *router.Context)
		DeleteCurrentUserEmailInfo() []oapi.Info
	}

	UserHandler interface {
		GetCurrentUser(c *router.Context)
		GetCurrentUserInfo() []oapi.Info

		PatchCurrentUser(c *router.Context)
		PatchCurrentUserInfo() []oapi.Info

		GetCurrentUserCredentials(c *router.Context)
		GetCurrentUserCredentialsInfo() []oapi.Info
	}
)
