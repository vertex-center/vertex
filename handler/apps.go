package handler

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/pkg/router"
)

type AppsHandler struct {
	appsService *service.AppsService
}

func NewAppsHandler(appsService *service.AppsService) port.AppsHandler {
	return &AppsHandler{
		appsService: appsService,
	}
}

func (h *AppsHandler) Get(c *router.Context) {
	c.JSON(h.appsService.All())
}
