package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type hardwareHandler struct {
	service port.HardwareService
}

func NewHardwareHandler(service port.HardwareService) port.HardwareHandler {
	return &hardwareHandler{
		service: service,
	}
}

func (h *hardwareHandler) GetHost() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) (*types.Host, error) {
		host, err := h.service.GetHost()
		if err != nil {
			return nil, err
		}
		return &host, nil
	})
}

func (h *hardwareHandler) GetCPUs() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) ([]types.CPU, error) {
		cpus, err := h.service.GetCPUs()
		if err != nil {
			return nil, err
		}
		return cpus, nil
	})
}

func (h *hardwareHandler) Reboot() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) error {
		return h.service.Reboot(ctx)
	})
}
