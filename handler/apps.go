package handler

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/pkg/router"
)

type appsHandler struct {
	appsService port.AppsService
}

func NewAppsHandler(appsService port.AppsService) port.AppsHandler {
	return &appsHandler{
		appsService: appsService,
	}
}

// docapi begin get_apps
// docapi method GET
// docapi summary Get all apps
// docapi tags Apps
// docapi response 200 {[]Meta} The list of apps.
// docapi end

func (h *appsHandler) GetApps(c *router.Context) {
	c.JSON(h.appsService.All())
}
