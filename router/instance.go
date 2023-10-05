package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addInstanceRoutes(r *router.Group) {
	r.GET("", handleGetInstance)
	r.DELETE("", handleDeleteInstance)
	r.PATCH("", handlePatchInstance)
	r.POST("/start", handleStartInstance)
	r.POST("/stop", handleStopInstance)
	r.PATCH("/environment", handlePatchEnvironment)
	r.GET("/events", headersSSE, handleInstanceEvents)
	r.GET("/docker", handleGetDocker)
	r.POST("/docker/recreate", handleRecreateDockerContainer)
	r.GET("/logs", handleGetLogs)
	r.POST("/update/service", handleUpdateService)
	r.GET("/versions", handleGetVersions)
}

// getParamInstanceUUID returns the UUID of the instance in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
func getParamInstanceUUID(c *router.Context) *uuid.UUID {
	p := c.Param("instance_uuid")
	if p == "" {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrInstanceUuidMissing,
			Message: "'instance_uuid' was missing in the URL",
		})
		return nil
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrInstanceUuidInvalid,
			Message: fmt.Sprintf("'%s' is not a valid UUID", p),
		})
		return nil
	}

	return &uid
}

// getInstance returns the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - instance_not_found: the instance with the given UUID was not found
//   - failed_to_retrieve_instance: failed to retrieve the instance from the database
func getInstance(c *router.Context) *types.Instance {
	instanceUUID := getParamInstanceUUID(c)
	if instanceUUID == nil {
		return nil
	}

	instance, err := instanceService.Get(*instanceUUID)
	if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
		_ = c.AbortWithError(http.StatusNotFound, api.Error{
			Code:    api.ErrInstanceNotFound,
			Message: fmt.Sprintf("instance %s not found", instanceUUID),
		})
		return nil
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetInstance,
			Message: fmt.Sprintf("failed to retrieve instance %s: %v", instanceUUID, err),
		})
		return nil
	}

	return instance
}

// handleGetInstance returns the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
func handleGetInstance(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}
	c.JSON(http.StatusOK, inst)
}

// handleDeleteInstance deletes the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - instance_not_found: the instance with the given UUID was not found
//   - instance_still_running: the instance with the given UUID is still running
//   - failed_to_delete_instance: failed to delete the instance from the database
func handleDeleteInstance(c *router.Context) {
	uid := getParamInstanceUUID(c)
	if uid == nil {
		return
	}

	err := instanceService.Delete(*uid)
	if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
		_ = c.AbortWithError(http.StatusNotFound, api.Error{
			Code:    api.ErrInstanceNotFound,
			Message: fmt.Sprintf("instance %s not found", uid),
		})
		return
	} else if err != nil && errors.Is(err, types.ErrInstanceStillRunning) {
		_ = c.AbortWithError(http.StatusConflict, api.Error{
			Code:    api.ErrInstanceStillRunning,
			Message: fmt.Sprintf("instance %s is still running", uid),
		})
		return
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToDeleteInstance,
			Message: fmt.Sprintf("failed to delete instance %s, %v", uid, err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

type handlePatchInstanceBody struct {
	LaunchOnStartup *bool                `json:"launch_on_startup,omitempty"`
	DisplayName     *string              `json:"display_name,omitempty"`
	Databases       map[string]uuid.UUID `json:"databases,omitempty"`
	Version         *string              `json:"version,omitempty"`
}

// handlePatchInstance updates the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - failed_to_parse_body: failed to parse the request body
//   - instance_not_found: the instance with the given UUID was not found
//   - failed_to_set_launch_on_startup: failed to set the launch on startup value
//   - failed_to_set_display_name: failed to set the display name
func handlePatchInstance(c *router.Context) {
	uid := getParamInstanceUUID(c)
	if uid == nil {
		return
	}

	inst := getInstance(c)
	if inst == nil {
		return
	}

	var body handlePatchInstanceBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	if body.LaunchOnStartup != nil {
		err = instanceSettingsService.SetLaunchOnStartup(inst, *body.LaunchOnStartup)
		if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
			_ = c.AbortWithError(http.StatusNotFound, api.Error{
				Code:    api.ErrInstanceNotFound,
				Message: fmt.Sprintf("instance %s not found", uid),
			})
			return
		} else if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
				Code:    api.ErrFailedToSetLaunchOnStartup,
				Message: fmt.Sprintf("failed to set launch on startup: %v", err),
			})
			return
		}
	}

	if body.DisplayName != nil {
		err = instanceSettingsService.SetDisplayName(inst, *body.DisplayName)
		if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
			_ = c.AbortWithError(http.StatusNotFound, api.Error{
				Code:    api.ErrInstanceNotFound,
				Message: fmt.Sprintf("instance %s not found", uid),
			})
			return
		} else if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
				Code:    api.ErrFailedToSetDisplayName,
				Message: fmt.Sprintf("failed to set display name: %v", err),
			})
			return
		}
	}

	if body.Databases != nil {
		err = instanceService.SetDatabases(inst, body.Databases)
		if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
			_ = c.AbortWithError(http.StatusNotFound, api.Error{
				Code:    api.ErrInstanceNotFound,
				Message: fmt.Sprintf("instance %s not found", uid),
			})
			return
		} else if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
				Code:    api.ErrFailedToSetDatabase,
				Message: fmt.Sprintf("failed to set databases: %v", err),
			})
			return
		}
	}

	if body.Version != nil {
		err = instanceSettingsService.SetVersion(inst, *body.Version)
		if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
			_ = c.AbortWithError(http.StatusNotFound, api.Error{
				Code:    api.ErrInstanceNotFound,
				Message: fmt.Sprintf("instance %s not found", uid),
			})
			return
		} else if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
				Code:    api.ErrFailedToSetVersion,
				Message: fmt.Sprintf("failed to set version: %v", err),
			})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// handleStartInstance starts the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - instance_not_found: the instance with the given UUID was not found
//   - instance_already_running: the instance with the given UUID is already running
//   - failed_to_start_instance: failed to start the instance
func handleStartInstance(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}

	err := instanceRunnerService.Start(inst)
	if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
		_ = c.AbortWithError(http.StatusNotFound, api.Error{
			Code:    api.ErrInstanceNotFound,
			Message: fmt.Sprintf("instance %s not found", inst.UUID),
		})
		return
	} else if err != nil && errors.Is(err, services.ErrInstanceAlreadyRunning) {
		_ = c.AbortWithError(http.StatusConflict, api.Error{
			Code:    api.ErrInstanceAlreadyRunning,
			Message: fmt.Sprintf("instance %s is already running", inst.UUID),
		})
		return
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToStartInstance,
			Message: fmt.Sprintf("failed to start instance %s: %v", inst.UUID, err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleStopInstance stops the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - instance_not_found: the instance with the given UUID was not found
//   - instance_not_running: the instance with the given UUID is not running
//   - failed_to_stop_instance: failed to stop the instance
func handleStopInstance(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}

	err := instanceRunnerService.Stop(inst)
	if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
		_ = c.AbortWithError(http.StatusNotFound, api.Error{
			Code:    api.ErrInstanceNotFound,
			Message: fmt.Sprintf("instance %s not found", inst.UUID),
		})
		return
	} else if err != nil && errors.Is(err, services.ErrInstanceNotRunning) {
		_ = c.AbortWithError(http.StatusConflict, api.Error{
			Code:    api.ErrInstanceNotRunning,
			Message: fmt.Sprintf("instance %s is not running", inst.UUID),
		})
		return
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToStopInstance,
			Message: fmt.Sprintf("failed to stop instance %s: %v", inst.UUID, err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handlePatchEnvironment updates the environment of the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - failed_to_parse_body: failed to parse the request body
//   - instance_not_found: the instance with the given UUID was not found
//   - failed_to_save_environment: failed to save the environment
//   - failed_to_recreate_container: failed to recreate the Docker container
func handlePatchEnvironment(c *router.Context) {
	var environment map[string]string
	err := c.BindJSON(&environment)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	inst := getInstance(c)
	if inst == nil {
		return
	}

	err = instanceEnvService.Save(inst, environment)
	if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
		_ = c.AbortWithError(http.StatusNotFound, api.Error{
			Code:    api.ErrInstanceNotFound,
			Message: fmt.Sprintf("instance %s not found", inst.UUID),
		})
		return
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToSetEnv,
			Message: fmt.Sprintf("failed to save environment: %v", err),
		})
		return
	}

	err = instanceRunnerService.RecreateContainer(inst)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToRecreateContainer,
			Message: fmt.Sprintf("failed to recreate container: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleInstanceEvents returns a stream of events for the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
func handleInstanceEvents(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}

	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := types.NewTempListener(func(e interface{}) {
		switch e := e.(type) {
		case types.EventInstanceLog:
			if inst.UUID != e.InstanceUUID {
				break
			}

			if e.Kind == types.LogKindOut || e.Kind == types.LogKindVertexOut {
				eventsChan <- sse.Event{
					Event: types.EventNameInstanceStdout,
					Data:  e.Message,
				}
			} else if e.Kind == types.LogKindErr || e.Kind == types.LogKindVertexErr {
				eventsChan <- sse.Event{
					Event: types.EventNameInstanceStderr,
					Data:  e.Message,
				}
			} else if e.Kind == types.LogKindDownload {
				eventsChan <- sse.Event{
					Event: types.EventNameInstanceDownload,
					Data:  e.Message,
				}
			}

		case types.EventInstanceStatusChange:
			if inst.UUID != e.InstanceUUID {
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

// handleGetDocker returns the Docker container info of the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - failed_to_get_docker_container_info: failed to get the Docker container info
func handleGetDocker(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}

	info, err := instanceRunnerService.GetDockerContainerInfo(*inst)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetContainerInfo,
			Message: fmt.Sprintf("failed to get docker container info: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// handleRecreateDockerContainer recreates the Docker container of the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - failed_to_recreate_container: failed to recreate the Docker container
func handleRecreateDockerContainer(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}

	err := instanceRunnerService.RecreateContainer(inst)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToRecreateContainer,
			Message: fmt.Sprintf("failed to recreate container: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleGetLogs returns the latest logs of the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - failed_to_get_logs: failed to get the logs
func handleGetLogs(c *router.Context) {
	uid := getParamInstanceUUID(c)
	if uid == nil {
		return
	}

	logs, err := instanceLogsService.GetLatestLogs(*uid)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetInstanceLogs,
			Message: fmt.Sprintf("failed to get logs: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// handleUpdateService updates the service of the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
func handleUpdateService(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}

	service, err := serviceService.GetById(inst.Service.ID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrServiceNotFound,
			Message: fmt.Sprintf("failed to get service: %v", err),
		})
		return
	}

	err = instanceServiceService.Update(inst, service)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToUpdateServiceInstance,
			Message: fmt.Sprintf("failed to update service: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleGetVersions returns the versions of the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
func handleGetVersions(c *router.Context) {
	inst := getInstance(c)
	if inst == nil {
		return
	}

	useCache := c.Query("reload") != "true"

	versions, err := instanceRunnerService.GetAllVersions(inst, useCache)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToGetVersions,
			Message: fmt.Sprintf("failed to get versions: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, versions)
}
