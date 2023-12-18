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

func (h *serviceHandler) GetService() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *GetServiceParams) (*types.Service, error) {
		return h.containerService.GetServiceByID(ctx, params.ServiceID)
	})
}

func (h *serviceHandler) GetServices() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) ([]types.Service, error) {
		return h.containerService.GetServices(ctx), nil
	})
}
