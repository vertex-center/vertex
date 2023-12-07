package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
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

func (h *appsHandler) GetApps(c *router.Context) {
	c.JSON(h.appsService.All())
}

func (h *appsHandler) GetAppsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get apps"),
		oapi.Description("Get all the apps installed on the server."),
		oapi.Response(http.StatusOK,
			oapi.WithResponseDesc("All apps installed on the server"),
			oapi.WithResponseModel([]app.Meta{}),
		),
	}
}
