package handler

import (
	"io"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
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

func (h *checksHandler) CheckInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("check"),
		fizz.Summary("Get all checks"),
		fizz.Description("Check that all vertex requirements are met."),
	}
}
