package handler

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
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

func (h *containersHandler) Get(c *router.Context) {
	installed := h.containerService.GetAll()
	c.JSON(installed)
}

func (h *containersHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get containers"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.Container{}),
		),
	}
}

func (h *containersHandler) GetTags(c *router.Context) {
	tags := h.containerService.GetTags()
	c.JSON(tags)
}

func (h *containersHandler) GetTagsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get tags"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]string{}),
		),
	}
}

func (h *containersHandler) Search(c *router.Context) {
	query := types.ContainerSearchQuery{}

	features := c.QueryArray("features[]")
	if len(features) > 0 {
		query.Features = &features
	}

	tags := c.QueryArray("tags[]")
	if len(tags) > 0 {
		query.Tags = &tags
	}

	installed := h.containerService.Search(query)
	c.JSON(installed)
}

func (h *containersHandler) SearchInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Search containers"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.Container{}),
		),
	}
}

func (h *containersHandler) CheckForUpdates(c *router.Context) {
	containers, err := h.containerService.CheckForUpdates()
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToCheckForUpdates,
			PublicMessage:  "Failed to check for updates.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(containers)
}

func (h *containersHandler) CheckForUpdatesInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Check for updates"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.Container{}),
		),
	}
}

func (h *containersHandler) Events(c *router.Context) {
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
}

func (h *containersHandler) EventsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get events"),
		oapi.Response(http.StatusOK),
	}
}
