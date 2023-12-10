package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/vertex/apps/logs/core/port"
)

type logsHandler struct {
	service port.LogsService
}

func NewLogsHandler(service port.LogsService) port.LogsHandler {
	return &logsHandler{
		service: service,
	}
}

type PushLogParams struct {
	Content string `json:"content" binding:"required"`
}

func (h *logsHandler) Push() gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, params *PushLogParams) error {
		return h.service.Push(params.Content)
	}, http.StatusOK)
}
