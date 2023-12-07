package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router/oapi"
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

func (h *appsHandler) GetAppsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getApps"),
		oapi.Summary("Get apps"),
		oapi.Description("Get all the apps installed on the server."),
	}
}
