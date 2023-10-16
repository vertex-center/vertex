package handler

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/pkg/router"
)

type AppsHandler struct {
	appsService port.AppsService
}

func NewAppsHandler(appsService port.AppsService) port.AppsHandler {
	return &AppsHandler{
		appsService: appsService,
	}
}

func (h *AppsHandler) Get(c *router.Context) {
	c.JSON(h.appsService.All())
}
