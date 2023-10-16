package handler

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	types2 "github.com/vertex-center/vertex/apps/containers/core/types"

	"github.com/vertex-center/vertex/pkg/router"
)

type ServiceHandler struct {
	serviceService   port.ServiceService
	containerService port.ContainerService
}

func NewServiceHandler(serviceService port.ServiceService, containerService port.ContainerService) port.ServiceHandler {
	return &ServiceHandler{
		serviceService:   serviceService,
		containerService: containerService,
	}
}

func (h *ServiceHandler) Get(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           types2.ErrCodeServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := h.serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types2.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(service)
}

func (h *ServiceHandler) Install(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           types2.ErrCodeServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := h.serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types2.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	inst, err := h.containerService.Install(service, "docker")
	if err != nil && errors.Is(err, types2.ErrServiceNotFound) {
		c.NotFound(router.Error{
			Code:           types2.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types2.ErrCodeFailedToInstallService,
			PublicMessage:  fmt.Sprintf("Failed to install service '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(inst)
}
