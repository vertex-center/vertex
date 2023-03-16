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
	"github.com/vertex-center/vertex/services/instances"
	servicesmanager "github.com/vertex-center/vertex/services/manager"
)

var logger = console.New("vertex::router")

func InitializeRouter() *gin.Engine {
	r := router.CreateRouter()
	r.Use(cors.Default())

	servicesGroup := r.Group("/services")
	servicesGroup.GET("", handleServicesInstalled)
	servicesGroup.GET("/available", handleServicesAvailable)
	servicesGroup.POST("/download", handleServiceDownload)

	serviceGroup := r.Group("/service/:service_uuid")
	serviceGroup.POST("/start", handleServiceStart)
	serviceGroup.POST("/stop", handleServiceStop)

	r.GET("/events", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", sse.ContentType)
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Origin", "*")
		c.Writer.Header().Set("X-Accel-Buffering", "no")
	}, handleEvents)

	return r
}

func handleServicesInstalled(c *gin.Context) {
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

	instance, err := instances.Install(body.Service)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "OK",
		"instance": instance,
	})
}

func handleServiceStart(c *gin.Context) {
	serviceUUIDParam := c.Param("service_uuid")
	if serviceUUIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("service_uuid was missing in the URL"))
		return
	}

	serviceUUID, err := uuid.Parse(serviceUUIDParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse service_uuid: %v", err))
		return
	}

	err = instances.Start(serviceUUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleServiceStop(c *gin.Context) {
	serviceUUIDParam := c.Param("service_uuid")
	if serviceUUIDParam == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("service_uuid was missing in the URL"))
		return
	}

	serviceUUID, err := uuid.Parse(serviceUUIDParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse service_uuid: %v", err))
		return
	}

	err = instances.Stop(serviceUUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleEvents(c *gin.Context) {
	instancesChan := make(chan instances.Event)
	id := instances.Register(instancesChan)

	done := c.Request.Context().Done()

	defer func() {
		instances.Unregister(id)
		close(instancesChan)
	}()

	c.Stream(func(w io.Writer) bool {
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
