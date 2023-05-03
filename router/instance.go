package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/types"
)

func addInstanceRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetInstance)
	r.DELETE("", handleDeleteInstance)
	r.PATCH("", handlePatchInstance)
	r.POST("/start", handleStartInstance)
	r.POST("/stop", handleStopInstance)
	r.PATCH("/environment", handlePatchEnvironment)
	r.GET("/events", headersSSE, handleInstanceEvents)
	r.GET("/dependencies", handleGetDependencies)
}

func getParamInstanceUUID(c *gin.Context) *uuid.UUID {
	p := c.Param("instance_uuid")
	if p == "" {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("instance_uuid was missing in the URL"))
		return nil
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse instance_uuid: %v", err))
		return nil
	}

	return &uid
}

func getInstance(c *gin.Context) *types.Instance {
	instanceUUID := getParamInstanceUUID(c)

	i, err := instanceService.Get(*instanceUUID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve instance %s: %v", instanceUUID, err))
		return nil
	}

	return i
}

func handleGetInstance(c *gin.Context) {
	i := getInstance(c)
	c.JSON(http.StatusOK, i)
}

func handleDeleteInstance(c *gin.Context) {
	uid := getParamInstanceUUID(c)

	err := instanceService.Delete(*uid)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete instance %s: %v", uid, err))
		return
	}

	c.Status(http.StatusOK)
}

type handlePatchInstanceBody struct {
	LaunchOnStartup *bool `json:"launch_on_startup"`
}

func handlePatchInstance(c *gin.Context) {
	uid := getParamInstanceUUID(c)

	var body handlePatchInstanceBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	if body.LaunchOnStartup != nil {
		err = instanceService.SetLaunchOnStartup(*uid, *body.LaunchOnStartup)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusOK)
}

func handleStartInstance(c *gin.Context) {
	uid := getParamInstanceUUID(c)

	err := instanceService.Start(*uid)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func handleStopInstance(c *gin.Context) {
	uid := getParamInstanceUUID(c)

	err := instanceService.Stop(*uid)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func handlePatchEnvironment(c *gin.Context) {
	var environment map[string]string
	err := c.BindJSON(&environment)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	uid := getParamInstanceUUID(c)

	err = instanceService.WriteEnv(*uid, environment)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save environment: %v", err))
		return
	}

	c.Status(http.StatusOK)
}

func handleInstanceEvents(c *gin.Context) {
	i := getInstance(c)

	instanceChan := make(chan types.InstanceEvent)
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
				logger.Error(err).Print()
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
				logger.Error(err).Print()
			}
			return true
		case <-done:
			return false
		}
	})
}

func handleGetDependencies(c *gin.Context) {
	i := getInstance(c)

	var deps = map[string]types.Package{}

	for name := range i.Dependencies {
		dep, err := packageService.Get(name)
		if err != nil {
			logger.Error(err).Print()
			continue
		}

		deps[name] = dep
	}

	c.JSON(http.StatusOK, deps)
}
