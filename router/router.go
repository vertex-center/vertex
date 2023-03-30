package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex/dependencies"
	"github.com/vertex-center/vertex/dependencies/dependency"
	"github.com/vertex-center/vertex/services/instance"
	"github.com/vertex-center/vertex/services/instances"
	servicesmanager "github.com/vertex-center/vertex/services/manager"
	"github.com/vertex-center/vertex/storage"
)

var logger = console.New("vertex::router")

func InitializeRouter() *gin.Engine {
	r, api := router.CreateRouter(cors.Default())
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.PathClient, "build"), true)))

	servicesGroup := api.Group("/services")
	servicesGroup.GET("/available", handleServicesAvailable)
	servicesGroup.POST("/download", handleServiceDownload)

	instancesGroup := api.Group("/instances")
	instancesGroup.GET("", handleGetInstances)
	instancesGroup.GET("/events", headersSSE, handleInstancesEvents)

	instanceGroup := api.Group("/instance/:instance_uuid")
	instanceGroup.GET("", handleGetInstance)
	instanceGroup.DELETE("", handleDeleteInstance)
	instanceGroup.POST("/start", handleStartInstance)
	instanceGroup.POST("/stop", handleStopInstance)
	instanceGroup.PATCH("/environment", handlePatchEnvironment)
	instanceGroup.GET("/events", headersSSE, handleInstanceEvents)
	instanceGroup.GET("/dependencies", handleGetDependencies)

	dependencyGroup := api.Group("/dependency/:dependency_id")
	dependencyGroup.POST("/install", handleInstallDependency)

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
	Repository string `json:"repository"`
}

func handleServiceDownload(c *gin.Context) {
	var body DownloadBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	i, err := instances.Install(body.Repository)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
		dep, err := dependencies.Get(name)
		if err != nil {
			logger.Error(err)
			continue
		}

		deps[name] = *dep
	}

	c.JSON(http.StatusOK, deps)
}

type InstallDependencyBody struct {
	PackageManager string `json:"package_manager"`
}

func handleInstallDependency(c *gin.Context) {
	var body InstallDependencyBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	dependencyID := c.Param("dependency_id")
	if dependencyID == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("dependency_id was missing in the URL"))
		return
	}

	dep, err := dependencies.Get(dependencyID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get dependency %s: %v", dependencyID, err))
		return
	}

	err = dep.Install(body.PackageManager)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
