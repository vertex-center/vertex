package router

import (
	types2 "github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	"io"

	"github.com/gin-contrib/sse"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

// handleGetContainers returns all installed containers.
func (r *AppRouter) handleGetContainers(c *router.Context) {
	installed := r.containerService.GetAll()
	c.JSON(installed)
}

func (r *AppRouter) handleGetTags(c *router.Context) {
	tags := r.containerService.GetTags()
	c.JSON(tags)
}

// handleSearchContainers returns all installed containers that match the query.
func (r *AppRouter) handleSearchContainers(c *router.Context) {
	query := types2.ContainerSearchQuery{}

	features := c.QueryArray("features[]")
	if len(features) > 0 {
		query.Features = &features
	}

	tags := c.QueryArray("tags[]")
	if len(tags) > 0 {
		query.Tags = &tags
	}

	installed := r.containerService.Search(query)
	c.JSON(installed)
}

// handleCheckForUpdates checks for updates for all installed containers.
// Errors can be:
//   - check_for_updates_failed
func (r *AppRouter) handleCheckForUpdates(c *router.Context) {
	containers, err := r.containerService.CheckForUpdates()
	if err != nil {
		c.Abort(router.Error{
			Code:           types2.ErrCodeFailedToCheckForUpdates,
			PublicMessage:  "Failed to check for updates.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(containers)
}

// handleContainersEvents returns a stream of events related to containers.
func (r *AppRouter) handleContainersEvents(c *router.Context) {
	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := vtypes.NewTempListener(func(e interface{}) {
		switch e.(type) {
		case types2.EventContainersChange:
			eventsChan <- sse.Event{
				Event: types2.EventNameContainersChange,
			}
		}
	})

	r.ctx.AddListener(listener)
	defer r.ctx.RemoveListener(listener)

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
