package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/services/instance"
	"github.com/vertex-center/vertex/services/instances"
	servicesmanager "github.com/vertex-center/vertex/services/manager"
)

var logger = console.New("vertex::router")

func InitializeRouter() *gin.Engine {
	r := router.CreateRouter()
	r.Use(cors.Default())

	servicesGroup := r.Group("/services")
	servicesGroup.GET("/available", handleServicesAvailable)
	servicesGroup.POST("/download", handleServiceDownload)

	instancesGroup := r.Group("/instances")
	instancesGroup.GET("", handleGetInstances)
	instancesGroup.GET("/events", headersSSE, handleInstancesEvents)

	serviceGroup := r.Group("/instance/:instance_uuid")
	serviceGroup.GET("", handleGetInstance)
	serviceGroup.DELETE("", handleDeleteInstance)
	serviceGroup.POST("/start", handleStartInstance)
	serviceGroup.POST("/stop", handleStopInstance)
	serviceGroup.GET("/events", headersSSE, handleInstanceEvents)

	return r
}

func headersSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

func handleGetInstances(c *gin.Context) {
	installed := instances.List()
	c.JSON(http.StatusOK, installed)
}

func handleServicesAvailable(c *gin.Context) {
	c.JSON(http.StatusOK, servicesmanager.ListAvailable())
}

type DownloadBody struct {
	Service services.Service `json:"service"`
}

func handleServiceDownload(c *gin.Context) {
	var body DownloadBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	i, err := instances.Install(body.Service)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "OK",
		"instance": i,
	})
}

func handleGetInstance(c *gin.Context) {
	instanceUUIDParam := c.Param("instance_uuid")
	if instanceUUIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("instance_uuid was missing in the URL"))
		return
	}

	instanceUUID, err := uuid.Parse(instanceUUIDParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse instance_uuid: %v", err))
		return
	}

	i, err := instances.Get(instanceUUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, i)
}

func handleDeleteInstance(c *gin.Context) {
	instanceUUIDParam := c.Param("instance_uuid")
	if instanceUUIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("instance_uuid was missing in the URL"))
		return
	}

	instanceUUID, err := uuid.Parse(instanceUUIDParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse instance_uuid: %v", err))
		return
	}

	err = instances.Delete(instanceUUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleStartInstance(c *gin.Context) {
	instanceUUIDParam := c.Param("instance_uuid")
	if instanceUUIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("instance_uuid was missing in the URL"))
		return
	}

	instanceUUID, err := uuid.Parse(instanceUUIDParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse instance_uuid: %v", err))
		return
	}

	err = instances.Start(instanceUUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleStopInstance(c *gin.Context) {
	instanceUUIDParam := c.Param("instance_uuid")
	if instanceUUIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("instance_uuid was missing in the URL"))
		return
	}

	instanceUUID, err := uuid.Parse(instanceUUIDParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse instance_uuid: %v", err))
		return
	}

	err = instances.Stop(instanceUUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleInstancesEvents(c *gin.Context) {
	instancesChan := make(chan instances.Event)
	id := instances.Register(instancesChan)

	done := c.Request.Context().Done()

	defer func() {
		instances.Unregister(id)
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

func handleInstanceEvents(c *gin.Context) {
	instanceUUIDParam := c.Param("instance_uuid")
	if instanceUUIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("instance_uuid was missing in the URL"))
		return
	}

	instanceUUID, err := uuid.Parse(instanceUUIDParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse instance_uuid: %v", err))
		return
	}

	i, err := instances.Get(instanceUUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get service %s: %v", instanceUUID, err))
		return
	}

	instanceChan := make(chan instance.Event)
	id := i.Register(instanceChan)

	defer func() {
		i.Unregister(id)
		close(instanceChan)
	}()

	done := c.Request.Context().Done()

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
		case e := <-instanceChan:
			err := sse.Encode(w, sse.Event{
				Event: e.Name,
				Data:  e.Data,
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
