package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type servicesHandler struct {
	containerService port.ContainerService
}

func NewServicesHandler(service port.ContainerService) port.ServicesHandler {
	return &servicesHandler{service}
}

func (h *servicesHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.Service, error) {
		return h.containerService.GetServices(c), nil
	})
}
