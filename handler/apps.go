package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/core/port"
	"github.com/wI2L/fizz"
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

func (h *appsHandler) GetAppsInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getApps"),
		fizz.Summary("Get apps"),
		fizz.Description("Get all the apps installed on the server."),
	}
}
