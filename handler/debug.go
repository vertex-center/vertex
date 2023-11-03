package handler

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/pkg/router"
)

type DebugHandler struct {
	debugService port.DebugService
}

func NewDebugHandler(debugService port.DebugService) port.DebugHandler {
	return &DebugHandler{
		debugService: debugService,
	}
}

// docapi begin hard_reset
// docapi method POST
// docapi summary Hard reset.
// docapi tags debug
// docapi response 200
// docapi end

func (h *DebugHandler) HardReset(c *router.Context) {
	h.debugService.HardReset()
	c.OK()
}
