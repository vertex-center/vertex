package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/core/port"
	"github.com/wI2L/fizz"
)

type debugHandler struct {
	debugService port.DebugService
}

func NewDebugHandler(debugService port.DebugService) port.DebugHandler {
	return &debugHandler{
		debugService: debugService,
	}
}

func (h *debugHandler) HardReset(c *gin.Context) error {
	h.debugService.HardReset()
	return nil
}

func (h *debugHandler) HardResetInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("hardReset"),
		fizz.Summary("Hard reset"),
		fizz.Description("This route allows deleting all the server data, which can be useful for debugging purposes."),
	}
}
