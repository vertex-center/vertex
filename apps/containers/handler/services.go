package handler

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/pkg/router"
)

type ServicesHandler struct {
	serviceService port.ServiceService
}

func NewServicesHandler(serviceService port.ServiceService) port.ServicesHandler {
	return &ServicesHandler{
		serviceService: serviceService,
	}
}

func (h *ServicesHandler) Get(c *router.Context) {
	c.JSON(h.serviceService.GetAll())
}
