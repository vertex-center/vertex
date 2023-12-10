package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/vertex/apps/logs/core/port"
)

type logsHandler struct {
	service  port.LogsService
	upgrader websocket.Upgrader
}

func NewLogsHandler(service port.LogsService) port.LogsHandler {
	return &logsHandler{
		service:  service,
		upgrader: websocket.Upgrader{},
	}
}

func (h *logsHandler) Push() gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return errors.Annotate(err, "upgrade connection")
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			err = h.service.Push(string(msg))
			if err != nil {
				return err
			}
		}
	}, http.StatusOK)
}
