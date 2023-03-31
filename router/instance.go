package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/dependency"
	"github.com/vertex-center/vertex/services/instance"
	"github.com/vertex-center/vertex/services/instances"
)

func addInstanceRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetInstance)
	r.DELETE("", handleDeleteInstance)
	r.POST("/start", handleStartInstance)
	r.POST("/stop", handleStopInstance)
	r.PATCH("/environment", handlePatchEnvironment)
	r.GET("/events", headersSSE, handleInstanceEvents)
	r.GET("/dependencies", handleGetDependencies)
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

	c.Status(http.StatusOK)
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

	c.Status(http.StatusOK)
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

	c.Status(http.StatusOK)
}

func handlePatchEnvironment(c *gin.Context) {
	var environment map[string]string
	err := c.BindJSON(&environment)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

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

	err = i.SetEnv(environment)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save environment: %v", err))
		return
	}

	c.Status(http.StatusOK)
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

func handleGetDependencies(c *gin.Context) {
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

	var deps = map[string]dependency.Dependency{}

	for name, _ := range i.Dependencies {
		dep, err := dependency.Get(name)
		if err != nil {
			logger.Error(err)
			continue
		}

		deps[name] = *dep
	}

	c.JSON(http.StatusOK, deps)
}
