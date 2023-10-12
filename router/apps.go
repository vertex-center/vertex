package router

import (
	"github.com/vertex-center/vertex/pkg/router"
)

func addAppsRoutes(r *router.Group) {
	r.GET("", getApps)
}

func getApps(c *router.Context) {
	c.JSON(appsService.All())
}
