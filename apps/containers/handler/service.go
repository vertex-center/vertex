package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
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

func (h *serviceHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get service"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Service{}),
		),
	}
}

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

func (h *serviceHandler) InstallInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Install service"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Container{}),
		),
	}
}
