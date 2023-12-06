package handler

import (
	"io"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"golang.org/x/net/context"
)

type ChecksHandler struct {
	checksService port.ChecksService
}

func NewChecksHandler(checksService port.ChecksService) port.ChecksHandler {
	return &ChecksHandler{
		checksService: checksService,
	}
}

// docapi begin admin_checks
// docapi method GET
// docapi summary Get all checks
// docapi desc Check that all vertex requirements are met.
// docapi tags Admin/Checks
// docapi response 200
// docapi end

func (h *ChecksHandler) Check(c *router.Context) {
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
}
