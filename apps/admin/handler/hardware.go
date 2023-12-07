package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type hardwareHandler struct {
	service port.HardwareService
}

func NewHardwareHandler(service port.HardwareService) port.HardwareHandler {
	return &hardwareHandler{
		service: service,
	}
}

func (h *hardwareHandler) GetHost(c *router.Context) {
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

func (h *hardwareHandler) GetHostInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get host"),
		oapi.Description("Get host information."),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Host{}),
		),
	}
}

func (h *hardwareHandler) GetCPUs(c *router.Context) {
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

func (h *hardwareHandler) GetCPUsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get CPUs"),
		oapi.Description("Get CPUs information."),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.CPU{}),
		),
	}
}

func (h *hardwareHandler) Reboot(c *router.Context) {
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

func (h *hardwareHandler) RebootInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Reboot"),
		oapi.Description("Reboot the host."),
		oapi.Response(http.StatusNoContent),
	}
}
