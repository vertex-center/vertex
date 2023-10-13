package router

import (
	"github.com/vertex-center/vertex/pkg/router"
)

// handleServicesAvailable handles the retrieval of all available services.
func (r *AppRouter) handleGetServices(c *router.Context) {
	c.JSON(r.serviceService.GetAll())
}
