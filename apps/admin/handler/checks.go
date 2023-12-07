package handler

import (
	"io"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
	"golang.org/x/net/context"
)

type checksHandler struct {
	checksService port.ChecksService
}

func NewChecksHandler(checksService port.ChecksService) port.ChecksHandler {
	return &checksHandler{
		checksService: checksService,
	}
}

func (h *checksHandler) Check() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) error {
		timeout, cancelTimeout := context.WithTimeout(c, 10*time.Second)
		resCh := h.checksService.CheckAll(timeout)
		defer cancelTimeout()

		c.Stream(func(w io.Writer) bool {
			res, ok := <-resCh
			if !ok {
				_ = sse.Encode(w, sse.Event{
					Event: "done",
				})
				return false
			}
			err := sse.Encode(w, sse.Event{
				Event: "check",
				Data:  res,
			})
			return err == nil
		})

		return nil
	})
}

func (h *checksHandler) CheckInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("check"),
		oapi.Summary("Get all checks"),
		oapi.Description("Check that all vertex requirements are met."),
	}
}
