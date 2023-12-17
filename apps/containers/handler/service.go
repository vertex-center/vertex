package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type serviceHandler struct {
	containerService port.ContainerService
}

func NewServiceHandler(containerService port.ContainerService) port.ServiceHandler {
	return &serviceHandler{containerService}
}

type GetServiceParams struct {
	ServiceID string `path:"service_id"`
}

func (h *serviceHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetServiceParams) (*types.Service, error) {
		return h.containerService.GetServiceByID(c, params.ServiceID)
	})
}

type InstallServiceParams struct {
	ServiceID string `path:"service_id"`
}

func (h *serviceHandler) Install() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InstallServiceParams) (*types.Container, error) {
		return h.containerService.Install(c, params.ServiceID)
	})
}
