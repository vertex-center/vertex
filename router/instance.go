package router

import (
	"errors"
	"fmt"
	"io"

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
		c.BadRequest(router.Error{
			Code:           api.ErrInstanceUuidMissing,
			PublicMessage:  "The request was missing the instance UUID.",
			PrivateMessage: "Field 'instance_uuid' is required.",
		})
		return nil
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           api.ErrInstanceUuidInvalid,
			PublicMessage:  "The instance UUID is invalid.",
			PrivateMessage: err.Error(),
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
		c.NotFound(router.Error{
			Code:           api.ErrInstanceNotFound,
			PublicMessage:  fmt.Sprintf("Instance %s not found.", instanceUUID),
			PrivateMessage: err.Error(),
		})
		return nil
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetInstance,
			PublicMessage:  fmt.Sprintf("Failed to retrieve instance %s.", instanceUUID),
			PrivateMessage: err.Error(),
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
	c.JSON(inst)
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
		c.NotFound(router.Error{
			Code:           api.ErrInstanceNotFound,
			PublicMessage:  fmt.Sprintf("Instance %s not found.", uid),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil && errors.Is(err, types.ErrInstanceStillRunning) {
		c.Conflict(router.Error{
			Code:           api.ErrInstanceStillRunning,
			PublicMessage:  fmt.Sprintf("Instance %s is still running.", uid),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToDeleteInstance,
			PublicMessage:  fmt.Sprintf("Failed to delete instance %s.", uid),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

type handlePatchInstanceBody struct {
	LaunchOnStartup *bool                `json:"launch_on_startup,omitempty"`
	DisplayName     *string              `json:"display_name,omitempty"`
	Databases       map[string]uuid.UUID `json:"databases,omitempty"`
	Version         *string              `json:"version,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
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
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	if body.LaunchOnStartup != nil {
		err = instanceSettingsService.SetLaunchOnStartup(inst, *body.LaunchOnStartup)
		if err != nil {
			c.Abort(router.Error{
				Code:           api.ErrFailedToSetLaunchOnStartup,
				PublicMessage:  "Failed to change launch on startup.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.DisplayName != nil {
		err = instanceSettingsService.SetDisplayName(inst, *body.DisplayName)
		if err != nil {
			c.Abort(router.Error{
				Code:           api.ErrFailedToSetDisplayName,
				PublicMessage:  "Failed to change display name.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Databases != nil {
		err = instanceService.SetDatabases(inst, body.Databases)
		if err != nil {
			c.Abort(router.Error{
				Code:           api.ErrFailedToSetDatabase,
				PublicMessage:  "Failed to change databases.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Version != nil {
		err = instanceSettingsService.SetVersion(inst, *body.Version)
		if err != nil {
			c.Abort(router.Error{
				Code:           api.ErrFailedToSetVersion,
				PublicMessage:  "Failed to change version.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Tags != nil {
		err = instanceSettingsService.SetTags(inst, body.Tags)
		if err != nil {
			c.Abort(router.Error{
				Code:           api.ErrFailedToSetTags,
				PublicMessage:  "Failed to change tags.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	c.OK()
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
		c.NotFound(router.Error{
			Code:           api.ErrInstanceNotFound,
			PublicMessage:  fmt.Sprintf("Instance %s not found.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil && errors.Is(err, services.ErrInstanceAlreadyRunning) {
		c.Conflict(router.Error{
			Code:           api.ErrInstanceAlreadyRunning,
			PublicMessage:  fmt.Sprintf("Instance %s is already running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToStartInstance,
			PublicMessage:  fmt.Sprintf("Failed to start instance %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
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
	if err != nil && errors.Is(err, services.ErrInstanceNotRunning) {
		c.Conflict(router.Error{
			Code:           api.ErrInstanceNotRunning,
			PublicMessage:  fmt.Sprintf("Instance %s is not running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToStopInstance,
			PublicMessage:  fmt.Sprintf("Failed to stop instance %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
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
	err := c.ParseBody(&environment)
	if err != nil {
		return
	}

	inst := getInstance(c)
	if inst == nil {
		return
	}

	err = instanceEnvService.Save(inst, environment)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToSetEnv,
			PublicMessage:  "failed to set environment",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = instanceRunnerService.RecreateContainer(inst)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToRecreateContainer,
			PublicMessage:  "Failed to recreate container.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
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
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetContainerInfo,
			PublicMessage:  fmt.Sprintf("Failed to get info for container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(info)
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
		c.Abort(router.Error{
			Code:           api.ErrFailedToRecreateContainer,
			PublicMessage:  fmt.Sprintf("Failed to recreate container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
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
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetInstanceLogs,
			PublicMessage:  fmt.Sprintf("Failed to get logs for instance %s.", uid),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(logs)
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
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service %s not found.", inst.Service.ID),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = instanceServiceService.Update(inst, service)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToUpdateServiceInstance,
			PublicMessage:  fmt.Sprintf("Failed to update service for instance %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
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
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetVersions,
			PublicMessage:  fmt.Sprintf("Failed to get versions for instance %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(versions)
}
