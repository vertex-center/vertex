package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
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
	r.GET("/docker", handleGetDocker)
	r.POST("/docker/recreate", handleRecreateDockerContainer)
	r.GET("/logs", handleGetLogs)
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
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

type handlePatchInstanceBody struct {
	LaunchOnStartup *bool   `json:"launch_on_startup,omitempty"`
	DisplayName     *string `json:"display_name,omitempty"`
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

	if body.DisplayName != nil {
		err = instanceService.SetDisplayName(*uid, *body.DisplayName)
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
	i := getInstance(c)

	err = instanceService.WriteEnv(*uid, environment)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save environment: %v", err))
		return
	}

	err = instanceService.RecreateContainer(i)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func handleInstanceEvents(c *gin.Context) {
	instance := getInstance(c)

	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := types.NewTempListener(func(e interface{}) {
		switch e := e.(type) {
		case types.EventInstanceLog:
			if instance.UUID != e.InstanceUUID {
				break
			}

			if e.Kind == types.LogKindOut || e.Kind == types.LogKindVertexOut {
				eventsChan <- sse.Event{
					Event: types.EventNameInstanceStdout,
					Data:  e.Message,
				}
			} else {
				eventsChan <- sse.Event{
					Event: types.EventNameInstanceStderr,
					Data:  e.Message,
				}
			}

		case types.EventInstanceStatusChange:
			if instance.UUID != e.InstanceUUID {
				break
			}

			eventsChan <- sse.Event{
				Event: types.EventNameInstanceStatusChange,
				Data:  e.Status,
			}
		}
	})

	eventInMemoryAdapter.AddListener(listener)
	defer eventInMemoryAdapter.RemoveListener(listener)

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

func handleGetDependencies(c *gin.Context) {
	i := getInstance(c)

	if i.Methods.Script == nil {
		_ = c.AbortWithError(http.StatusNotFound, errors.New("this service doesn't use scripts, so it doesn't have dependencies"))
	}

	var deps = map[string]types.Package{}

	if i.Methods.Script.Dependencies != nil {
		for name := range *i.Methods.Script.Dependencies {
			dep, err := packageService.GetByID(name)
			if err != nil {
				log.Error(err)
				continue
			}
			deps[name] = dep
		}
	}

	c.JSON(http.StatusOK, deps)
}

func handleGetDocker(c *gin.Context) {
	uid := getParamInstanceUUID(c)

	info, err := instanceService.GetDockerContainerInfo(*uid)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

func handleRecreateDockerContainer(c *gin.Context) {
	i := getInstance(c)

	err := instanceService.RecreateContainer(i)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func handleGetLogs(c *gin.Context) {
	uid := getParamInstanceUUID(c)

	logs, err := instanceService.GetLatestLogs(*uid)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, logs)
}
