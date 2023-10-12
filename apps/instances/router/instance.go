package router

import (
	"errors"
	"fmt"
	"io"

	"github.com/gin-contrib/sse"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/service"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	vtypes "github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

// getParamInstanceUUID returns the UUID of the instance in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
func (r *AppRouter) getParamInstanceUUID(c *router.Context) *uuid.UUID {
	p := c.Param("instance_uuid")
	if p == "" {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeInstanceUuidMissing,
			PublicMessage:  "The request was missing the instance UUID.",
			PrivateMessage: "Field 'instance_uuid' is required.",
		})
		return nil
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeInstanceUuidInvalid,
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
func (r *AppRouter) getInstance(c *router.Context) *types.Instance {
	instanceUUID := r.getParamInstanceUUID(c)
	if instanceUUID == nil {
		return nil
	}

	instance, err := r.instanceService.Get(*instanceUUID)
	if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
		c.NotFound(router.Error{
			Code:           types.ErrCodeInstanceNotFound,
			PublicMessage:  fmt.Sprintf("The instance '%s' could not be found.", instanceUUID),
			PrivateMessage: err.Error(),
		})
		return nil
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetInstance,
			PublicMessage:  fmt.Sprintf("Failed to retrieve instance '%s'.", instanceUUID),
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
func (r *AppRouter) handleGetInstance(c *router.Context) {
	inst := r.getInstance(c)
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
func (r *AppRouter) handleDeleteInstance(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	err := r.instanceService.Delete(inst)
	if err != nil && errors.Is(err, types.ErrInstanceStillRunning) {
		c.Conflict(router.Error{
			Code:           types.ErrCodeInstanceStillRunning,
			PublicMessage:  fmt.Sprintf("The instance '%s' is still running. Stop it first before deleting.", inst.DisplayName),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToDeleteInstance,
			PublicMessage:  fmt.Sprintf("The instance '%s' could not be deleted.", inst.DisplayName),
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
func (r *AppRouter) handlePatchInstance(c *router.Context) {
	uid := r.getParamInstanceUUID(c)
	if uid == nil {
		return
	}

	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	var body handlePatchInstanceBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	if body.LaunchOnStartup != nil {
		err = r.instanceSettingsService.SetLaunchOnStartup(inst, *body.LaunchOnStartup)
		if err != nil {
			c.Abort(router.Error{
				Code:           types.ErrCodeFailedToSetLaunchOnStartup,
				PublicMessage:  "Failed to change launch on startup.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.DisplayName != nil {
		err = r.instanceSettingsService.SetDisplayName(inst, *body.DisplayName)
		if err != nil {
			c.Abort(router.Error{
				Code:           types.ErrCodeFailedToSetDisplayName,
				PublicMessage:  "Failed to change display name.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Databases != nil {
		err = r.instanceService.SetDatabases(inst, body.Databases)
		if err != nil {
			c.Abort(router.Error{
				Code:           types.ErrCodeFailedToSetDatabase,
				PublicMessage:  "Failed to change databases.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Version != nil {
		err = r.instanceSettingsService.SetVersion(inst, *body.Version)
		if err != nil {
			c.Abort(router.Error{
				Code:           types.ErrCodeFailedToSetVersion,
				PublicMessage:  "Failed to change version.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.Tags != nil {
		err = r.instanceSettingsService.SetTags(inst, body.Tags)
		if err != nil {
			c.Abort(router.Error{
				Code:           types.ErrCodeFailedToSetTags,
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
func (r *AppRouter) handleStartInstance(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	err := r.instanceRunnerService.Start(inst)
	if err != nil && errors.Is(err, types.ErrInstanceNotFound) {
		c.NotFound(router.Error{
			Code:           types.ErrCodeInstanceNotFound,
			PublicMessage:  fmt.Sprintf("Instance '%s' not found.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil && errors.Is(err, service.ErrInstanceAlreadyRunning) {
		c.Conflict(router.Error{
			Code:           types.ErrCodeInstanceAlreadyRunning,
			PublicMessage:  fmt.Sprintf("Instance %s is already running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToStartInstance,
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
func (r *AppRouter) handleStopInstance(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	err := r.instanceRunnerService.Stop(inst)
	if err != nil && errors.Is(err, service.ErrInstanceNotRunning) {
		c.Conflict(router.Error{
			Code:           types.ErrCodeInstanceNotRunning,
			PublicMessage:  fmt.Sprintf("Instance %s is not running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToStopInstance,
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
func (r *AppRouter) handlePatchEnvironment(c *router.Context) {
	var environment map[string]string
	err := c.ParseBody(&environment)
	if err != nil {
		return
	}

	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	err = r.instanceEnvService.Save(inst, environment)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToSetEnv,
			PublicMessage:  "failed to set environment",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = r.instanceRunnerService.RecreateContainer(inst)
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
func (r *AppRouter) handleInstanceEvents(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := vtypes.NewTempListener(func(e interface{}) {
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

// handleGetDocker returns the Docker container info of the instance with the UUID in the URL.
// Errors can be:
//   - missing_instance_uuid: the instance_uuid parameter was missing in the URL
//   - invalid_instance_uuid: the instance_uuid parameter was not a valid UUID
//   - failed_to_get_docker_container_info: failed to get the Docker container info
func (r *AppRouter) handleGetDocker(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	info, err := r.instanceRunnerService.GetDockerContainerInfo(*inst)
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
func (r *AppRouter) handleRecreateDockerContainer(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	err := r.instanceRunnerService.RecreateContainer(inst)
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
func (r *AppRouter) handleGetLogs(c *router.Context) {
	uid := r.getParamInstanceUUID(c)
	if uid == nil {
		return
	}

	logs, err := r.instanceLogsService.GetLatestLogs(*uid)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetInstanceLogs,
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
func (r *AppRouter) handleUpdateService(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	serv, err := r.serviceService.GetById(inst.Service.ID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service %s not found.", inst.Service.ID),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = r.instanceServiceService.Update(inst, serv)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToUpdateServiceInstance,
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
func (r *AppRouter) handleGetVersions(c *router.Context) {
	inst := r.getInstance(c)
	if inst == nil {
		return
	}

	useCache := c.Query("reload") != "true"

	versions, err := r.instanceRunnerService.GetAllVersions(inst, useCache)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetVersions,
			PublicMessage:  fmt.Sprintf("Failed to get versions for instance %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(versions)
}
