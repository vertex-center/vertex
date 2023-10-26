package handler

import (
	"io"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"

	"github.com/gin-contrib/sse"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

type ContainersHandler struct {
	ctx              *apptypes.Context
	containerService port.ContainerService
}

func NewContainersHandler(ctx *apptypes.Context, containerService port.ContainerService) port.ContainersHandler {
	return &ContainersHandler{
		ctx:              ctx,
		containerService: containerService,
	}
}

func (h *ContainersHandler) Get(c *router.Context) {
	installed := h.containerService.GetAll()
	c.JSON(installed)
}

func (h *ContainersHandler) GetTags(c *router.Context) {
	tags := h.containerService.GetTags()
	c.JSON(tags)
}

func (h *ContainersHandler) Search(c *router.Context) {
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

func (h *ContainersHandler) CheckForUpdates(c *router.Context) {
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

func (h *ContainersHandler) Events(c *router.Context) {
	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := vtypes.NewTempListener(func(e interface{}) {
		switch e.(type) {
		case types.EventContainersChange:
			eventsChan <- sse.Event{
				Event: types.EventNameContainersChange,
			}
		}
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
