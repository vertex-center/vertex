package router

import (
	"errors"
	"fmt"
	"github.com/vertex-center/vertex/apps/containers/core/service"
	types3 "github.com/vertex-center/vertex/apps/containers/core/types"
	types2 "github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"io"

	"github.com/gin-contrib/sse"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

// getParamContainerUUID returns the UUID of the container in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
func (r *AppRouter) getParamContainerUUID(c *router.Context) *uuid.UUID {
	p := c.Param("container_uuid")
	if p == "" {
		c.BadRequest(router.Error{
			Code:           types3.ErrCodeContainerUuidMissing,
			PublicMessage:  "The request was missing the container UUID.",
			PrivateMessage: "Field 'container_uuid' is required.",
		})
		return nil
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types3.ErrCodeContainerUuidInvalid,
			PublicMessage:  "The container UUID is invalid.",
			PrivateMessage: err.Error(),
		})
		return nil
	}

	return &uid
}

// getContainer returns the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - container_not_found: the container with the given UUID was not found
//   - failed_to_retrieve_container: failed to retrieve the container from the database
func (r *AppRouter) getContainer(c *router.Context) *types3.Container {
	containerUUID := r.getParamContainerUUID(c)
	if containerUUID == nil {
		return nil
	}

	container, err := r.containerService.Get(*containerUUID)
	if err != nil && errors.Is(err, types3.ErrContainerNotFound) {
		c.NotFound(router.Error{
			Code:           types3.ErrCodeContainerNotFound,
			PublicMessage:  fmt.Sprintf("The container '%s' could not be found.", containerUUID),
			PrivateMessage: err.Error(),
		})
		return nil
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToGetContainer,
			PublicMessage:  fmt.Sprintf("Failed to retrieve container '%s'.", containerUUID),
			PrivateMessage: err.Error(),
		})
		return nil
	}

	return container
}

// handleGetContainer returns the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
func (r *AppRouter) handleGetContainer(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}
	c.JSON(inst)
}

// handleDeleteContainer deletes the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - container_not_found: the container with the given UUID was not found
//   - container_still_running: the container with the given UUID is still running
//   - failed_to_delete_container: failed to delete the container from the database
func (r *AppRouter) handleDeleteContainer(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	err := r.containerService.Delete(inst)
	if err != nil && errors.Is(err, types3.ErrContainerStillRunning) {
		c.Conflict(router.Error{
			Code:           types3.ErrCodeContainerStillRunning,
			PublicMessage:  fmt.Sprintf("The container '%s' is still running. Stop it first before deleting.", inst.DisplayName),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToDeleteContainer,
			PublicMessage:  fmt.Sprintf("The container '%s' could not be deleted.", inst.DisplayName),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

type handlePatchContainerBody struct {
	LaunchOnStartup *bool                `json:"launch_on_startup,omitempty"`
	DisplayName     *string              `json:"display_name,omitempty"`
	Databases       map[string]uuid.UUID `json:"databases,omitempty"`
	Version         *string              `json:"version,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
}

// handlePatchContainer updates the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - failed_to_parse_body: failed to parse the request body
//   - container_not_found: the container with the given UUID was not found
//   - failed_to_set_launch_on_startup: failed to set the launch on startup value
//   - failed_to_set_display_name: failed to set the display name
func (r *AppRouter) handlePatchContainer(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	var body handlePatchContainerBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	if body.LaunchOnStartup != nil {
		err = r.containerSettingsService.SetLaunchOnStartup(inst, *body.LaunchOnStartup)
		if err != nil {
			c.Abort(router.Error{
				Code:           types3.ErrCodeFailedToSetLaunchOnStartup,
				PublicMessage:  "Failed to change launch on startup.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.DisplayName != nil && *body.DisplayName != "" {
		err = r.containerSettingsService.SetDisplayName(inst, *body.DisplayName)
		if err != nil {
			c.Abort(router.Error{
				Code:           types3.ErrCodeFailedToSetDisplayName,
				PublicMessage:  "Failed to change display name.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Databases != nil {
		err = r.containerService.SetDatabases(inst, body.Databases)
		if err != nil {
			c.Abort(router.Error{
				Code:           types3.ErrCodeFailedToSetDatabase,
				PublicMessage:  "Failed to change databases.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Version != nil {
		err = r.containerSettingsService.SetVersion(inst, *body.Version)
		if err != nil {
			c.Abort(router.Error{
				Code:           types3.ErrCodeFailedToSetVersion,
				PublicMessage:  "Failed to change version.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Tags != nil {
		err = r.containerSettingsService.SetTags(inst, body.Tags)
		if err != nil {
			c.Abort(router.Error{
				Code:           types3.ErrCodeFailedToSetTags,
				PublicMessage:  "Failed to change tags.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	c.OK()
}

// handleStartContainer starts the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - container_not_found: the container with the given UUID was not found
//   - container_already_running: the container with the given UUID is already running
//   - failed_to_start_container: failed to start the container
func (r *AppRouter) handleStartContainer(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	err := r.containerRunnerService.Start(inst)
	if err != nil && errors.Is(err, types3.ErrContainerNotFound) {
		c.NotFound(router.Error{
			Code:           types3.ErrCodeContainerNotFound,
			PublicMessage:  fmt.Sprintf("Container '%s' not found.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil && errors.Is(err, service.ErrContainerAlreadyRunning) {
		c.Conflict(router.Error{
			Code:           types3.ErrCodeContainerAlreadyRunning,
			PublicMessage:  fmt.Sprintf("Container %s is already running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToStartContainer,
			PublicMessage:  fmt.Sprintf("Failed to start container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// handleStopContainer stops the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - container_not_found: the container with the given UUID was not found
//   - container_not_running: the container with the given UUID is not running
//   - failed_to_stop_container: failed to stop the container
func (r *AppRouter) handleStopContainer(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	err := r.containerRunnerService.Stop(inst)
	if err != nil && errors.Is(err, service.ErrContainerNotRunning) {
		c.Conflict(router.Error{
			Code:           types3.ErrCodeContainerNotRunning,
			PublicMessage:  fmt.Sprintf("Container %s is not running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToStopContainer,
			PublicMessage:  fmt.Sprintf("Failed to stop container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// handlePatchEnvironment updates the environment of the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - failed_to_parse_body: failed to parse the request body
//   - container_not_found: the container with the given UUID was not found
//   - failed_to_save_environment: failed to save the environment
//   - failed_to_recreate_container: failed to recreate the Docker container
func (r *AppRouter) handlePatchEnvironment(c *router.Context) {
	var environment map[string]string
	err := c.ParseBody(&environment)
	if err != nil {
		return
	}

	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	err = r.containerEnvService.Save(inst, environment)
	if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToSetEnv,
			PublicMessage:  "failed to set environment",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = r.containerRunnerService.RecreateContainer(inst)
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

// handleContainerEvents returns a stream of events for the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
func (r *AppRouter) handleContainerEvents(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := types2.NewTempListener(func(e interface{}) {
		switch e := e.(type) {
		case types3.EventContainerLog:
			if inst.UUID != e.ContainerUUID {
				break
			}

			if e.Kind == types3.LogKindOut || e.Kind == types3.LogKindVertexOut {
				eventsChan <- sse.Event{
					Event: types3.EventNameContainerStdout,
					Data:  e.Message,
				}
			} else if e.Kind == types3.LogKindErr || e.Kind == types3.LogKindVertexErr {
				eventsChan <- sse.Event{
					Event: types3.EventNameContainerStderr,
					Data:  e.Message,
				}
			} else if e.Kind == types3.LogKindDownload {
				eventsChan <- sse.Event{
					Event: types3.EventNameContainerDownload,
					Data:  e.Message,
				}
			}

		case types3.EventContainerStatusChange:
			if inst.UUID != e.ContainerUUID {
				break
			}

			eventsChan <- sse.Event{
				Event: types3.EventNameContainerStatusChange,
				Data:  e.Status,
			}
		}
	})

	r.ctx.AddListener(listener)
	defer r.ctx.RemoveListener(listener)

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

// handleGetDocker returns the Docker container info of the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - failed_to_get_docker_container_info: failed to get the Docker container info
func (r *AppRouter) handleGetDocker(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	info, err := r.containerRunnerService.GetDockerContainerInfo(*inst)
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

// handleRecreateDockerContainer recreates the Docker container of the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - failed_to_recreate_container: failed to recreate the Docker container
func (r *AppRouter) handleRecreateDockerContainer(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	err := r.containerRunnerService.RecreateContainer(inst)
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

// handleGetLogs returns the latest logs of the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
//   - invalid_container_uuid: the container_uuid parameter was not a valid UUID
//   - failed_to_get_logs: failed to get the logs
func (r *AppRouter) handleGetLogs(c *router.Context) {
	uid := r.getParamContainerUUID(c)
	if uid == nil {
		return
	}

	logs, err := r.containerLogsService.GetLatestLogs(*uid)
	if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToGetContainerLogs,
			PublicMessage:  fmt.Sprintf("Failed to get logs for container %s.", uid),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(logs)
}

// handleUpdateService updates the service of the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
func (r *AppRouter) handleUpdateService(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	serv, err := r.serviceService.GetById(inst.Service.ID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types3.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service %s not found.", inst.Service.ID),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = r.containerServiceService.Update(inst, serv)
	if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToUpdateServiceContainer,
			PublicMessage:  fmt.Sprintf("Failed to update service for container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// handleGetVersions returns the versions of the container with the UUID in the URL.
// Errors can be:
//   - missing_container_uuid: the container_uuid parameter was missing in the URL
func (r *AppRouter) handleGetVersions(c *router.Context) {
	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	useCache := c.Query("reload") != "true"

	versions, err := r.containerRunnerService.GetAllVersions(inst, useCache)
	if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToGetVersions,
			PublicMessage:  fmt.Sprintf("Failed to get versions for container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(versions)
}

func (r *AppRouter) handleWaitContainer(c *router.Context) {
	cond := c.Param("cond")

	inst := r.getContainer(c)
	if inst == nil {
		return
	}

	err := r.containerRunnerService.WaitCondition(inst, types2.WaitContainerCondition(cond))
	if err != nil {
		c.Abort(router.Error{
			Code:           types3.ErrCodeFailedToWaitContainer,
			PublicMessage:  fmt.Sprintf("Failed to wait the event '%s' for container %s.", cond, inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
