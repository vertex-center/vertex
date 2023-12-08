package handler

import (
	"io"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
)

type containerHandler struct {
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
	return &containerHandler{
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

type GetContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetContainerParams) (*types.Container, error) {
		return h.containerService.Get(*params.ContainerUUID)
	})
}

func (h *containerHandler) GetInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getContainer"),
		fizz.Summary("Get a container"),
	}
}

type DeleteContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) Delete() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteContainerParams) error {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}
		return h.containerService.Delete(inst)
	})
}

func (h *containerHandler) DeleteInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("deleteContainer"),
		fizz.Summary("Delete a container"),
	}
}

// User can also add alternate username,password
type PatchBodyDatabase struct {
	ContainerID  uuid.UUID `json:"container_id"`
	DatabaseName *string   `json:"db_name"`
}

type PatchContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`

	LaunchOnStartup *bool                        `json:"launch_on_startup,omitempty"`
	DisplayName     *string                      `json:"display_name,omitempty"`
	Databases       map[string]PatchBodyDatabase `json:"databases,omitempty"`
	Version         *string                      `json:"version,omitempty"`
	Tags            []string                     `json:"tags,omitempty"`
}

func (h *containerHandler) Patch() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchContainerParams) error {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}

		if params.LaunchOnStartup != nil {
			err = h.containerSettingsService.SetLaunchOnStartup(inst, *params.LaunchOnStartup)
			if err != nil {
				return err
			}
		}

		if params.DisplayName != nil && *params.DisplayName != "" {
			err = h.containerSettingsService.SetDisplayName(inst, *params.DisplayName)
			if err != nil {
				return err
			}
		}

		if params.Databases != nil {
			databases := map[string]uuid.UUID{}
			options := map[string]*types.SetDatabasesOptions{}

			for databaseID, container := range params.Databases {
				databases[databaseID] = container.ContainerID
				options[databaseID] = &types.SetDatabasesOptions{
					DatabaseName: container.DatabaseName,
				}
			}

			err = h.containerService.SetDatabases(inst, databases, options)
			if err != nil {
				return err
			}
		}

		if params.Version != nil {
			err = h.containerSettingsService.SetVersion(inst, *params.Version)
			if err != nil {
				return err
			}
		}

		if params.Tags != nil {
			err = h.containerSettingsService.SetTags(inst, params.Tags)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (h *containerHandler) PatchInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("patchContainer"),
		fizz.Summary("Patch a container"),
	}
}

type StartContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) Start() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StartContainerParams) error {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}
		return h.containerRunnerService.Start(inst)
	})
}

func (h *containerHandler) StartInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("startContainer"),
		fizz.Summary("Start a container"),
	}
}

type StopContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) Stop() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StopContainerParams) error {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}
		return h.containerRunnerService.Stop(inst)
	})
}

func (h *containerHandler) StopInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("stopContainer"),
		fizz.Summary("Stop a container"),
	}
}

type PatchEnvironmentBody map[string]string

func (h *containerHandler) PatchEnvironment() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchContainerParams) error {
		var environment PatchEnvironmentBody
		err := c.BindJSON(&environment)
		if err != nil {
			return err
		}

		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}

		err = h.containerEnvService.Save(inst, types.ContainerEnvVariables(environment))
		if err != nil {
			return err
		}

		err = h.containerRunnerService.RecreateContainer(inst)
		if err != nil {
			return err
		}
		return nil
	})
}

func (h *containerHandler) PatchEnvironmentInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("patchContainerEnvironment"),
		fizz.Summary("Patch a container environment"),
	}
}

type EventsContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) Events() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *EventsContainerParams) error {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}

		eventsChan := make(chan sse.Event)
		defer close(eventsChan)

		done := c.Request.Context().Done()

		listener := event.NewTempListener(func(e event.Event) error {
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
			return nil
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

		return nil
	})
}

func (h *containerHandler) EventsInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("eventsContainer"),
		fizz.Summary("Get container events"),
		fizz.Description("Get events for a container, sent as Server-Sent Events (SSE)."),
	}
}

type DockerContainerInfo map[string]any

func (h *containerHandler) GetDocker() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetContainerParams) (map[string]any, error) {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return nil, err
		}

		return h.containerRunnerService.GetDockerContainerInfo(*inst)
	})
}

func (h *containerHandler) GetDockerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getDockerContainer"),
		fizz.Summary("Get Docker container info"),
	}
}

type RecreateContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) RecreateDocker() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *RecreateContainerParams) error {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}
		return h.containerRunnerService.RecreateContainer(inst)
	})
}

func (h *containerHandler) RecreateDockerInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("recreateDockerContainer"),
		fizz.Summary("Recreate Docker container"),
	}
}

type LogsContainerParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) GetLogs() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *LogsContainerParams) ([]types.LogLine, error) {
		return h.containerLogsService.GetLatestLogs(*params.ContainerUUID)
	})
}

func (h *containerHandler) GetLogsInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getContainerLogs"),
		fizz.Summary("Get container logs"),
	}
}

type UpdateServiceParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) UpdateService() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *UpdateServiceParams) error {
		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}

		serv, err := h.serviceService.GetById(inst.Service.ID)
		if err != nil {
			return err
		}

		return h.containerServiceService.Update(inst, serv)
	})
}

func (h *containerHandler) UpdateServiceInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("updateService"),
		fizz.Summary("Update service"),
	}
}

type GetVersionsParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) GetVersions() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetVersionsParams) ([]string, error) {
		useCache := c.Query("reload") != "true"

		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return nil, err
		}
		return h.containerRunnerService.GetAllVersions(inst, useCache)
	})
}

func (h *containerHandler) GetVersionsInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getContainerVersions"),
		fizz.Summary("Get container versions"),
	}
}

type WaitStatusParams struct {
	ContainerUUID *uuid.UUID `path:"container_uuid"`
}

func (h *containerHandler) WaitStatus() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *WaitStatusParams) error {
		status := c.Query("status")

		inst, err := h.containerService.Get(*params.ContainerUUID)
		if err != nil {
			return err
		}

		return h.containerRunnerService.WaitStatus(inst, status)
	})
}

func (h *containerHandler) WaitStatusInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("waitContainerStatus"),
		fizz.Summary("Wait for a status change"),
	}
}
