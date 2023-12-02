package handler

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type HardwareHandler struct {
	hardwareService port.HardwareService
}

func NewHardwareHandler(hardwareService port.HardwareService) port.HardwareHandler {
	return &HardwareHandler{
		hardwareService: hardwareService,
	}
}

// docapi begin get_host
// docapi method GET
// docapi summary Get host
// docapi tags Apps/Admin/Hardware
// docapi response 200 {Host} The host information.
// docapi response 500
// docapi end

func (h *HardwareHandler) GetHost(c *router.Context) {
	host, err := h.hardwareService.GetHost()
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
// docapi tags Apps/Admin/Hardware
// docapi response 200 {[]CPU} The CPUs information.
// docapi response 500
// docapi end

func (h *HardwareHandler) GetCPUs(c *router.Context) {
	cpus, err := h.hardwareService.GetCPUs()
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
