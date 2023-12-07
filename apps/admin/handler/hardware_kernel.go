package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
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

func (h *hardwareKernelHandler) RebootInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("reboot"),
		oapi.Summary("Reboot"),
		oapi.Description("Reboot the host."),
	}
}
