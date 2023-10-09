package router

import (
	"github.com/vertex-center/vertex/pkg/router"
)

func addServicesRoutes(r *router.Group) {
	r.GET("", handleGetServices)
	r.Static("/icons", "./live/services/icons")
}

// handleServicesAvailable handles the retrieval of all available services.
func handleGetServices(c *router.Context) {
	c.JSON(serviceService.GetAll())
}
