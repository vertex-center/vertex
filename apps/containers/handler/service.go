package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	apierrors "github.com/juju/errors"
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

type GetServiceParams struct {
	ServiceID string `path:"service_id"`
}

func (h *serviceHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetServiceParams) (*types.Service, error) {
		service, err := h.serviceService.GetById(params.ServiceID)
		if err != nil {
			return nil, apierrors.NewNotFound(err, "service not found")
		}
		return &service, nil
	})
}

func (h *serviceHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getService"),
		oapi.Summary("Get service"),
	}
}

type InstallServiceParams struct {
	ServiceID string `path:"service_id"`
}

func (h *serviceHandler) Install() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InstallServiceParams) (*types.Container, error) {
		service, err := h.serviceService.GetById(params.ServiceID)
		if err != nil {
			return nil, apierrors.NewNotFound(err, "service not found")
		}

		inst, err := h.containerService.Install(service, "docker")
		if err != nil && errors.Is(err, types.ErrServiceNotFound) {
			return nil, apierrors.NewNotFound(err, "service not found")
		} else if err != nil {
			return nil, apierrors.Annotate(err, "failed to install service")
		}
		return inst, nil
	})
}

func (h *serviceHandler) InstallInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("installService"),
		oapi.Summary("Install service"),
	}
}
