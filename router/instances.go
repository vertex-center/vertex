package router

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/services/instance"
)

func addInstancesRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetInstances)
	r.GET("/events", headersSSE, handleInstancesEvents)
}

func handleGetInstances(c *gin.Context) {
	installed := instance.All()
	c.JSON(http.StatusOK, installed)
}

func handleInstancesEvents(c *gin.Context) {
	instancesChan := make(chan instance.Event)
	id := instance.Register(instancesChan)

	done := c.Request.Context().Done()

	defer func() {
		instance.Unregister(id)
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
