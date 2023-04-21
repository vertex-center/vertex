package router

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/types"
)

func addInstancesRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetInstances)
	r.GET("/events", headersSSE, handleInstancesEvents)
}

func handleGetInstances(c *gin.Context) {
	installed := instanceService.GetAll()
	c.JSON(http.StatusOK, installed)
}

func handleInstancesEvents(c *gin.Context) {
	instancesChan := make(chan types.InstanceEvent)
	id := instanceService.AddListener(instancesChan)

	done := c.Request.Context().Done()

	defer func() {
		instanceService.RemoveListener(id)
		close(instancesChan)
	}()

	first := true

	c.Stream(func(w io.Writer) bool {
		if first {
			err := sse.Encode(w, sse.Event{
				Event: "open",
			})

			if err != nil {
				logger.Error(err)
				return false
			}
			first = false
			return true
		}

		select {
		case e := <-instancesChan:
			err := sse.Encode(w, sse.Event{
				Event: e.Name,
				Data:  e.Name,
			})
			if err != nil {
				logger.Error(err)
			}
			return true
		case <-done:
			return false
		}
	})
}
