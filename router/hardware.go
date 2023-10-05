package router

import (
	"github.com/vertex-center/vertex/pkg/router"
)

func addHardwareRoutes(r *router.Group) {
	r.GET("", handleGetHardware)
}

func handleGetHardware(c *router.Context) {
	c.JSON(hardwareService.GetHardware())
}
