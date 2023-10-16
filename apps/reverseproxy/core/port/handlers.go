package port

import (
	"github.com/vertex-center/vertex/pkg/router"
)

type (
	ProxyHandler interface {
		GetRedirects(c *router.Context)
		AddRedirect(c *router.Context)
		RemoveRedirect(c *router.Context)
	}
)
