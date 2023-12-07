package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/core/types/api"
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

func (h *hardwareKernelHandler) RebootInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Reboot"),
		oapi.Description("Reboot the host."),
		oapi.Response(http.StatusNoContent),
	}
}
