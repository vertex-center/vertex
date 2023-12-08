package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
)

type servicesHandler struct {
	serviceService port.ServiceService
}

func NewServicesHandler(serviceService port.ServiceService) port.ServicesHandler {
	return &servicesHandler{
		serviceService: serviceService,
	}
}

func (h *servicesHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.Service, error) {
		return h.serviceService.GetAll(), nil
	})
}

func (h *servicesHandler) GetInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getServices"),
		fizz.Summary("Get services"),
	}
}
