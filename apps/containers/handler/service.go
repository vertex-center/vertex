package handler

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type serviceHandler struct {
	serviceService   port.ServiceService
	containerService port.ContainerService
}

func NewServiceHandler(serviceService port.ServiceService, containerService port.ContainerService) port.ServiceHandler {
	return &serviceHandler{
		serviceService:   serviceService,
		containerService: containerService,
	}
}

// docapi begin vx_containers_get_service
// docapi method GET
// docapi summary Get service
// docapi tags Containers
// docapi query service_id {string} The service ID.
// docapi response 200 {Service} The service.
// docapi response 400
// docapi response 404
// docapi end

func (h *serviceHandler) Get(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := h.serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(service)
}

// docapi begin vx_containers_install_service
// docapi method POST
// docapi summary Install a service
// docapi tags Containers
// docapi query service_id {string} The service ID.
// docapi response 200 {Container} The container.
// docapi response 400
// docapi response 404
// docapi response 500
// docapi end

func (h *serviceHandler) Install(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := h.serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	inst, err := h.containerService.Install(service, "docker")
	if err != nil && errors.Is(err, types.ErrServiceNotFound) {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToInstallService,
			PublicMessage:  fmt.Sprintf("Failed to install service '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(inst)
}
