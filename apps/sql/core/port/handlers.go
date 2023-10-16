package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	DBMSHandler interface {
		Get(c *router.Context)
		Install(c *router.Context)
	}
)
