package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	MetricsHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		InstallCollector(c *router.Context)
		InstallCollectorInfo() []oapi.Info

		InstallVisualizer(c *router.Context)
		InstallVisualizerInfo() []oapi.Info
	}
)
