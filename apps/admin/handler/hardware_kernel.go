package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
)

type hardwareKernelHandler struct {
	service port.HardwareKernelService
}

func NewHardwareKernelHandler(service port.HardwareKernelService) port.HardwareKernelHandler {
	return &hardwareKernelHandler{
		service: service,
	}
}

func (h *hardwareKernelHandler) Reboot() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) error {
		return h.service.Reboot()
	})
}
