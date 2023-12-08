package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
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

func (h *hardwareKernelHandler) RebootInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("reboot"),
		fizz.Summary("Reboot"),
		fizz.Description("Reboot the host."),
	}
}
