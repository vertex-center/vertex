package handler

import (
	"io"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/router"
)

type containersHandler struct {
	ctx              *apptypes.Context
	containerService port.ContainerService
}

func NewContainersHandler(ctx *apptypes.Context, containerService port.ContainerService) port.ContainersHandler {
	return &containersHandler{
		ctx:              ctx,
		containerService: containerService,
	}
}

func (h *containersHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (types.Containers, error) {
		return h.containerService.GetContainers(c)
	})
}

func (h *containersHandler) Search() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (types.Containers, error) {
		query := types.ContainerSearchQuery{}

		features := c.QueryArray("features[]")
		if len(features) > 0 {
			query.Features = &features
		}

		tags := c.QueryArray("tags[]")
		if len(tags) > 0 {
			query.Tags = &tags
		}

		return h.containerService.Search(c, query)
	})
}

func (h *containersHandler) CheckForUpdates() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (types.Containers, error) {
		return h.containerService.CheckForUpdates(c)
	})
}

func (h *containersHandler) Events() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) error {
		eventsChan := make(chan sse.Event)
		defer close(eventsChan)

		done := c.Request.Context().Done()

		listener := event.NewTempListener(func(e event.Event) error {
			switch e.(type) {
			case types.EventContainersChange:
				eventsChan <- sse.Event{
					Event: types.EventNameContainersChange,
				}
			}
			return nil
		})

		h.ctx.AddListener(listener)
		defer h.ctx.RemoveListener(listener)

		first := true

		c.Stream(func(w io.Writer) bool {
			if first {
				err := sse.Encode(w, sse.Event{
					Event: "open",
				})

				if err != nil {
					log.Error(err)
					return false
				}
				first = false
				return true
			}

			select {
			case e := <-eventsChan:
				err := sse.Encode(w, e)
				if err != nil {
					log.Error(err)
				}
				return true
			case <-done:
				return false
			}
		})

		return nil
	})
}
