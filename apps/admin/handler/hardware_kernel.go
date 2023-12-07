package handler

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/core/types/api"
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

// docapi begin reboot_kernel
// docapi method POST
// docapi summary Reboot
// docapi tags Hardware
// docapi response 204
// docapi response 500
// docapi end

func (h *hardwareKernelHandler) Reboot(c *router.Context) {
	err := h.service.Reboot()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToReboot,
			PublicMessage:  "Failed to reboot.",
			PrivateMessage: err.Error(),
		})
		return
	}
	c.OK()
}
