package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	DBMSHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		Install(c *router.Context)
		InstallInfo() []oapi.Info
	}
)
