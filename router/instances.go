package router

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/types"
)

func addInstancesRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetInstances)
	r.GET("/checkupdates", handleCheckForUpdates)
	r.GET("/events", headersSSE, handleInstancesEvents)
}

func handleGetInstances(c *gin.Context) {
	installed := instanceService.GetAll()
	c.JSON(http.StatusOK, installed)
}

func handleCheckForUpdates(c *gin.Context) {
	instances, err := instanceService.CheckForUpdates()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, instances)
}

func handleInstancesEvents(c *gin.Context) {
	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := types.NewTempListener(func(e interface{}) {
		switch e.(type) {
		case types.EventInstancesChange:
			eventsChan <- sse.Event{
				Event: types.EventNameInstancesChange,
			}
		}
	})

	eventInMemoryRepo.AddListener(listener)
	defer eventInMemoryRepo.RemoveListener(listener)

	first := true

	c.Stream(func(w io.Writer) bool {
		if first {
			err := sse.Encode(w, sse.Event{
				Event: "open",
			})

			if err != nil {
				logger.Error(err).Print()
				return false
			}
			first = false
			return true
		}

		select {
		case e := <-eventsChan:
			err := sse.Encode(w, e)
			if err != nil {
				logger.Error(err).Print()
			}
			return true
		case <-done:
			return false
		}
	})
}
