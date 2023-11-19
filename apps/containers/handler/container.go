package handler

import (
	"errors"
	"fmt"
	"io"

	"github.com/gin-contrib/sse"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/service"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

type ContainerHandler struct {
	ctx                      *apptypes.Context
	containerService         port.ContainerService
	containerSettingsService port.ContainerSettingsService
	containerRunnerService   port.ContainerRunnerService
	containerEnvService      port.ContainerEnvService
	containerServiceService  port.ContainerServiceService
	containerLogsService     port.ContainerLogsService
	serviceService           port.ServiceService
}

type ContainerHandlerParams struct {
	Ctx                      *apptypes.Context
	ContainerService         port.ContainerService
	ContainerSettingsService port.ContainerSettingsService
	ContainerRunnerService   port.ContainerRunnerService
	ContainerEnvService      port.ContainerEnvService
	ContainerServiceService  port.ContainerServiceService
	ContainerLogsService     port.ContainerLogsService
	ServiceService           port.ServiceService
}

func NewContainerHandler(params ContainerHandlerParams) port.ContainerHandler {
	return &ContainerHandler{
		ctx:                      params.Ctx,
		containerService:         params.ContainerService,
		containerSettingsService: params.ContainerSettingsService,
		containerRunnerService:   params.ContainerRunnerService,
		containerEnvService:      params.ContainerEnvService,
		containerServiceService:  params.ContainerServiceService,
		containerLogsService:     params.ContainerLogsService,
		serviceService:           params.ServiceService,
	}
}

func (h *ContainerHandler) getParamContainerUUID(c *router.Context) *uuid.UUID {
	p := c.Param("container_uuid")
	if p == "" {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeContainerUuidMissing,
			PublicMessage:  "The request was missing the container UUID.",
			PrivateMessage: "Field 'container_uuid' is required.",
		})
		return nil
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeContainerUuidInvalid,
			PublicMessage:  "The container UUID is invalid.",
			PrivateMessage: err.Error(),
		})
		return nil
	}

	return &uid
}

func (h *ContainerHandler) getContainer(c *router.Context) *types.Container {
	containerUUID := h.getParamContainerUUID(c)
	if containerUUID == nil {
		return nil
	}

	container, err := h.containerService.Get(*containerUUID)
	if err != nil && errors.Is(err, types.ErrContainerNotFound) {
		c.NotFound(router.Error{
			Code:           types.ErrCodeContainerNotFound,
			PublicMessage:  fmt.Sprintf("The container '%s' could not be found.", containerUUID),
			PrivateMessage: err.Error(),
		})
		return nil
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetContainer,
			PublicMessage:  fmt.Sprintf("Failed to retrieve container '%s'.", containerUUID),
			PrivateMessage: err.Error(),
		})
		return nil
	}

	return container
}

// docapi begin vx_containers_get_container
// docapi method GET
// docapi summary Get a container
// docapi tags Apps/Containers
// docapi response 200 {Container} The container.
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) Get(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}
	c.JSON(inst)
}

// docapi begin vx_containers_delete_container
// docapi method DELETE
// docapi summary Delete a container
// docapi tags Apps/Containers
// docapi response 204
// docapi response 404
// docapi response 409 {Error} The container is still running.
// docapi response 500
// docapi end

func (h *ContainerHandler) Delete(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	err := h.containerService.Delete(inst)
	if err != nil && errors.Is(err, types.ErrContainerStillRunning) {
		c.Conflict(router.Error{
			Code:           types.ErrCodeContainerStillRunning,
			PublicMessage:  fmt.Sprintf("The container '%s' is still running. Stop it first before deleting.", inst.DisplayName),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToDeleteContainer,
			PublicMessage:  fmt.Sprintf("The container '%s' could not be deleted.", inst.DisplayName),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

type PatchBody struct {
	LaunchOnStartup *bool                        `json:"launch_on_startup,omitempty"`
	DisplayName     *string                      `json:"display_name,omitempty"`
	Databases       map[string]PatchBodyDatabase `json:"databases,omitempty"`
	Version         *string                      `json:"version,omitempty"`
	Tags            []string                     `json:"tags,omitempty"`
}

// User can also add alternate username,password
type PatchBodyDatabase struct {
	ContainerID  uuid.UUID `json:"container_id"`
	DatabaseName *string   `json:"db_name"`
}

// docapi begin vx_containers_patch_container
// docapi method PATCH
// docapi summary Patch a container
// docapi tags Apps/Containers
// docapi body {PatchBody} The container patch.
// docapi response 204
// docapi response 400
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) Patch(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	var body PatchBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	if body.LaunchOnStartup != nil {
		err = h.containerSettingsService.SetLaunchOnStartup(inst, *body.LaunchOnStartup)
		if err != nil {
			c.Abort(router.Error{
				Code:           types.ErrCodeFailedToSetLaunchOnStartup,
				PublicMessage:  "Failed to change launch on startup.",
				PrivateMessage: err.Error(),
			})
			return
		}
	}

	if body.DisplayName != nil && *body.DisplayName != "" {
		err = h.containerSettingsService.SetDisplayName(inst, *body.DisplayName)
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

		databases := map[string]uuid.UUID{}
		options := map[string]*types.SetDatabasesOptions{}

		for databaseID, container := range body.Databases {
			databases[databaseID] = container.ContainerID
			options[databaseID] = &types.SetDatabasesOptions{
				DatabaseName: container.DatabaseName,
			}
		}

		err = h.containerService.SetDatabases(inst, databases, options)
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
		err = h.containerSettingsService.SetVersion(inst, *body.Version)
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
		err = h.containerSettingsService.SetTags(inst, body.Tags)
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

// docapi begin vx_containers_start_container
// docapi method POST
// docapi summary Start a container
// docapi tags Apps/Containers
// docapi response 204
// docapi response 404
// docapi response 409 {Error} The container is already running.
// docapi response 500
// docapi end

func (h *ContainerHandler) Start(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	err := h.containerRunnerService.Start(inst)
	if err != nil && errors.Is(err, types.ErrContainerNotFound) {
		c.NotFound(router.Error{
			Code:           types.ErrCodeContainerNotFound,
			PublicMessage:  fmt.Sprintf("Container '%s' not found.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil && errors.Is(err, service.ErrContainerAlreadyRunning) {
		c.Conflict(router.Error{
			Code:           types.ErrCodeContainerAlreadyRunning,
			PublicMessage:  fmt.Sprintf("Container %s is already running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToStartContainer,
			PublicMessage:  fmt.Sprintf("Failed to start container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// docapi begin vx_containers_stop_container
// docapi method POST
// docapi summary Stop a container
// docapi tags Apps/Containers
// docapi response 204
// docapi response 404
// docapi response 409 {Error} The container is not running.
// docapi response 500
// docapi end

func (h *ContainerHandler) Stop(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	err := h.containerRunnerService.Stop(inst)
	if err != nil && errors.Is(err, service.ErrContainerNotRunning) {
		c.Conflict(router.Error{
			Code:           types.ErrCodeContainerNotRunning,
			PublicMessage:  fmt.Sprintf("Container %s is not running.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToStopContainer,
			PublicMessage:  fmt.Sprintf("Failed to stop container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

type PatchEnvironmentBody map[string]string

// docapi begin vx_containers_patch_environment
// docapi method PATCH
// docapi summary Patch a container environment
// docapi tags Apps/Containers
// docapi body {PatchEnvironmentBody} The environment variables to set.
// docapi response 204
// docapi response 400
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) PatchEnvironment(c *router.Context) {
	var environment PatchEnvironmentBody
	err := c.ParseBody(&environment)
	if err != nil {
		return
	}

	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	err = h.containerEnvService.Save(inst, types.ContainerEnvVariables(environment))
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToSetEnv,
			PublicMessage:  "failed to set environment",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = h.containerRunnerService.RecreateContainer(inst)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToRecreateContainer,
			PublicMessage:  "Failed to recreate container.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// docapi begin vx_containers_events_container
// docapi method GET
// docapi summary Get container events
// docapi desc Get events for a container, sent as Server-Sent Events (SSE).
// docapi tags Apps/Containers
// docapi response 200
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) Events(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	eventsChan := make(chan sse.Event)
	defer close(eventsChan)

	done := c.Request.Context().Done()

	listener := event.NewTempListener(func(e event.Event) {
		switch e := e.(type) {
		case types.EventContainerLog:
			if inst.UUID != e.ContainerUUID {
				break
			}

			if e.Kind == types.LogKindOut || e.Kind == types.LogKindVertexOut {
				eventsChan <- sse.Event{
					Event: types.EventNameContainerStdout,
					Data:  e.Message,
				}
			} else if e.Kind == types.LogKindErr || e.Kind == types.LogKindVertexErr {
				eventsChan <- sse.Event{
					Event: types.EventNameContainerStderr,
					Data:  e.Message,
				}
			} else if e.Kind == types.LogKindDownload {
				eventsChan <- sse.Event{
					Event: types.EventNameContainerDownload,
					Data:  e.Message,
				}
			}

		case types.EventContainerStatusChange:
			if inst.UUID != e.ContainerUUID {
				break
			}

			eventsChan <- sse.Event{
				Event: types.EventNameContainerStatusChange,
				Data:  e.Status,
			}
		}
	})

	h.ctx.AddListener(listener)
	defer h.ctx.RemoveListener(listener)

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

type DockerContainerInfo map[string]any

// docapi begin vx_containers_get_docker
// docapi method GET
// docapi summary Get Docker container info
// docapi tags Apps/Containers
// docapi response 200 {DockerContainerInfo} The Docker container info.
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) GetDocker(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	info, err := h.containerRunnerService.GetDockerContainerInfo(*inst)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetContainerInfo,
			PublicMessage:  fmt.Sprintf("Failed to get info for container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(info)
}

// docapi begin vx_containers_recreate_docker
// docapi method POST
// docapi summary Recreate Docker container
// docapi tags Apps/Containers
// docapi response 204
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) RecreateDocker(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	err := h.containerRunnerService.RecreateContainer(inst)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToRecreateContainer,
			PublicMessage:  fmt.Sprintf("Failed to recreate container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// docapi begin vx_containers_get_logs
// docapi method GET
// docapi summary Get container logs
// docapi tags Apps/Containers
// docapi response 200 {[]LogLine} The logs.
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) GetLogs(c *router.Context) {
	uid := h.getParamContainerUUID(c)
	if uid == nil {
		return
	}

	logs, err := h.containerLogsService.GetLatestLogs(*uid)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetContainerLogs,
			PublicMessage:  fmt.Sprintf("Failed to get logs for container %s.", uid),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(logs)
}

// docapi begin vx_containers_update_service
// docapi method POST
// docapi summary Update service
// docapi tags Apps/Containers
// docapi response 204
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) UpdateService(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	serv, err := h.serviceService.GetById(inst.Service.ID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service %s not found.", inst.Service.ID),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = h.containerServiceService.Update(inst, serv)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToUpdateServiceContainer,
			PublicMessage:  fmt.Sprintf("Failed to update service for container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// docapi begin vx_containers_get_versions
// docapi method GET
// docapi summary Get container versions
// docapi tags Apps/Containers
// docapi query reload {bool} Whether to reload the versions from the registry or use the cache.
// docapi response 200 {[]string} The versions.
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) GetVersions(c *router.Context) {
	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	useCache := c.Query("reload") != "true"

	versions, err := h.containerRunnerService.GetAllVersions(inst, useCache)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetVersions,
			PublicMessage:  fmt.Sprintf("Failed to get versions for container %s.", inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(versions)
}

// docapi begin vx_containers_wait
// docapi method GET
// docapi summary Wait for a container event
// docapi tags Apps/Containers
// docapi query cond {string} The condition to wait for.
// docapi response 204
// docapi response 404
// docapi response 500
// docapi end

func (h *ContainerHandler) Wait(c *router.Context) {
	cond := c.Param("cond")

	inst := h.getContainer(c)
	if inst == nil {
		return
	}

	err := h.containerRunnerService.WaitCondition(inst, types.WaitContainerCondition(cond))
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToWaitContainer,
			PublicMessage:  fmt.Sprintf("Failed to wait the event '%s' for container %s.", cond, inst.UUID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}