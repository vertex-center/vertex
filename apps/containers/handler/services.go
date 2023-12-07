package handler

import (
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/pkg/router"
)

type servicesHandler struct {
	serviceService port.ServiceService
}

func NewServicesHandler(serviceService port.ServiceService) port.ServicesHandler {
	return &servicesHandler{
		serviceService: serviceService,
	}
}

// docapi begin vx_containers_get_services
// docapi method GET
// docapi summary Get services
// docapi tags Containers
// docapi response 200 {[]Service} The services.
// docapi end

func (h *servicesHandler) Get(c *router.Context) {
	c.JSON(h.serviceService.GetAll())
}
