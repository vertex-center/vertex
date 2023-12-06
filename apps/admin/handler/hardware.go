package handler

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type HardwareHandler struct {
	service port.HardwareService
}

func NewHardwareHandler(service port.HardwareService) port.HardwareHandler {
	return &HardwareHandler{
		service: service,
	}
}

// docapi begin get_host
// docapi method GET
// docapi summary Get host
// docapi tags Hardware
// docapi response 200 {Host} The host information.
// docapi response 500
// docapi end

func (h *HardwareHandler) GetHost(c *router.Context) {
	host, err := h.service.GetHost()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetHost,
			PublicMessage:  "Failed to get host information.",
			PrivateMessage: err.Error(),
		})
		return
	}
	c.JSON(host)
}

// docapi begin get_cpus
// docapi method GET
// docapi summary Get CPUs
// docapi tags Hardware
// docapi response 200 {[]CPU} The CPUs information.
// docapi response 500
// docapi end

func (h *HardwareHandler) GetCPUs(c *router.Context) {
	cpus, err := h.service.GetCPUs()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetCPUs,
			PublicMessage:  "Failed to get CPUs information.",
			PrivateMessage: err.Error(),
		})
		return
	}
	c.JSON(cpus)
}

// docapi begin reboot
// docapi method POST
// docapi summary Reboot
// docapi tags Hardware
// docapi response 204
// docapi response 500
// docapi end

func (h *HardwareHandler) Reboot(c *router.Context) {
	err := h.service.Reboot(c)
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
