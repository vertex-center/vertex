package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/core/port"
)

type appsHandler struct {
	appsService port.AppsService
}

func NewAppsHandler(appsService port.AppsService) port.AppsHandler {
	return &appsHandler{
		appsService: appsService,
	}
}

func (h *appsHandler) GetApps(*gin.Context) ([]app.Meta, error) {
	return h.appsService.All(), nil
}
