package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type servicesHandler struct {
	serviceService port.ServiceService
}

func NewServicesHandler(serviceService port.ServiceService) port.ServicesHandler {
	return &servicesHandler{
		serviceService: serviceService,
	}
}

func (h *servicesHandler) Get(c *router.Context) {
	c.JSON(h.serviceService.GetAll())
}

func (h *servicesHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get services"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.Service{}),
		),
	}
}
