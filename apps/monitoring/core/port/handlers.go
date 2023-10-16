package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	MetricsHandler interface {
		Get(c *router.Context)
		InstallCollector(c *router.Context)
		InstallVisualizer(c *router.Context)
	}
)
