package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	ProxyHandler interface {
		GetRedirects(c *router.Context)
		GetRedirectsInfo() []oapi.Info

		AddRedirect(c *router.Context)
		AddRedirectInfo() []oapi.Info

		RemoveRedirect(c *router.Context)
		RemoveRedirectInfo() []oapi.Info
	}
)
