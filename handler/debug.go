package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type debugHandler struct {
	debugService port.DebugService
}

func NewDebugHandler(debugService port.DebugService) port.DebugHandler {
	return &debugHandler{
		debugService: debugService,
	}
}

func (h *debugHandler) HardReset(c *router.Context) {
	h.debugService.HardReset()
	c.OK()
}

func (h *debugHandler) HardResetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Hard reset"),
		oapi.Description("This route allows deleting all the server data, which can be useful for debugging purposes."),
		oapi.Response(http.StatusNoContent),
	}
}
