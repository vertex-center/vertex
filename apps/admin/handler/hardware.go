package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
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

func (h *hardwareHandler) GetHost() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*types.Host, error) {
		host, err := h.service.GetHost()
		if err != nil {
			return nil, err
		}
		return &host, nil
	})
}

func (h *hardwareHandler) GetHostInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getHost"),
		oapi.Summary("Get host"),
		oapi.Description("Get host information."),
	}
}

func (h *hardwareHandler) GetCPUs() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.CPU, error) {
		cpus, err := h.service.GetCPUs()
		if err != nil {
			return nil, err
		}
		return cpus, nil
	})
}

func (h *hardwareHandler) GetCPUsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getCPUs"),
		oapi.Summary("Get CPUs"),
		oapi.Description("Get CPUs information."),
	}
}

func (h *hardwareHandler) Reboot() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) error {
		return h.service.Reboot(c)
	})
}

func (h *hardwareHandler) RebootInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("reboot"),
		oapi.Summary("Reboot"),
		oapi.Description("Reboot the host."),
	}
}
