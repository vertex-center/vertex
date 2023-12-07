package handler

import (
	"io"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
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
	return router.Handler(func(c *gin.Context) (map[uuid.UUID]*types.Container, error) {
		return h.containerService.GetAll(), nil
	})
}

func (h *containersHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getContainers"),
		oapi.Summary("Get containers"),
	}
}

func (h *containersHandler) GetTags() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]string, error) {
		return h.containerService.GetTags(), nil
	})
}

func (h *containersHandler) GetTagsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getTags"),
		oapi.Summary("Get tags"),
	}
}

func (h *containersHandler) Search() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (map[uuid.UUID]*types.Container, error) {
		query := types.ContainerSearchQuery{}

		features := c.QueryArray("features[]")
		if len(features) > 0 {
			query.Features = &features
		}

		tags := c.QueryArray("tags[]")
		if len(tags) > 0 {
			query.Tags = &tags
		}

		return h.containerService.Search(query), nil
	})
}

func (h *containersHandler) SearchInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("searchContainers"),
		oapi.Summary("Search containers"),
	}
}

func (h *containersHandler) CheckForUpdates() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (map[uuid.UUID]*types.Container, error) {
		return h.containerService.CheckForUpdates()
	})
}

func (h *containersHandler) CheckForUpdatesInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("checkForUpdates"),
		oapi.Summary("Check for updates"),
	}
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

func (h *containersHandler) EventsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("events"),
		oapi.Summary("Get events"),
	}
}
