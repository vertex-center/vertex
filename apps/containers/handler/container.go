package handler

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/router"
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
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) Get() gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, params *GetContainerParams) (*types.Container, error) {
		return h.containerService.Get(c, params.ContainerUUID.UUID)
	}, http.StatusOK)
}

type DeleteContainerParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) Delete() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return err
		}
		return h.containerService.Delete(c, inst)
	})
}

// User can also add alternate username,password
type PatchBodyDatabase struct {
	ContainerID  uuid.UUID `json:"container_id"`
	DatabaseName *string   `json:"db_name"`
}

type PatchContainerParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`

	LaunchOnStartup *bool                        `json:"launch_on_startup,omitempty"`
	DisplayName     *string                      `json:"display_name,omitempty"`
	Databases       map[string]PatchBodyDatabase `json:"databases,omitempty"`
	Version         *string                      `json:"version,omitempty"`
	Tags            []string                     `json:"tags,omitempty"`
}

func (h *containerHandler) Patch() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
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

			err = h.containerService.SetDatabases(c, inst, databases, options)
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

type StartContainerParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) Start() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StartContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return err
		}
		return h.containerRunnerService.Start(c, inst)
	})
}

type StopContainerParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) Stop() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StopContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return err
		}
		return h.containerRunnerService.Stop(c, inst)
	})
}

type PatchEnvironmentParams struct {
	ContainerUUID uuid.NullUUID               `path:"container_uuid"`
	Env           types.ContainerEnvVariables `body:"env"`
}

func (h *containerHandler) PatchEnvironment() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchEnvironmentParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return err
		}

		err = h.containerEnvService.Save(inst, params.Env)
		if err != nil {
			return err
		}

		err = h.containerRunnerService.RecreateContainer(c, inst)
		if err != nil {
			return err
		}
		return nil
	})
}

type EventsContainerParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) Events() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *EventsContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
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

func (h *containerHandler) GetDocker() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetContainerParams) (map[string]any, error) {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return nil, err
		}

		return h.containerRunnerService.GetDockerContainerInfo(c, *inst)
	})
}

type RecreateContainerParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) RecreateDocker() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *RecreateContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return err
		}
		return h.containerRunnerService.RecreateContainer(c, inst)
	})
}

type LogsContainerParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) GetLogs() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *LogsContainerParams) ([]types.LogLine, error) {
		return h.containerLogsService.GetLatestLogs(params.ContainerUUID.UUID)
	})
}

type UpdateServiceParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) UpdateService() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *UpdateServiceParams) error {
		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
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

type GetVersionsParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) GetVersions() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetVersionsParams) ([]string, error) {
		useCache := c.Query("reload") != "true"

		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return nil, err
		}
		return h.containerRunnerService.GetAllVersions(c, inst, useCache)
	})
}

type WaitStatusParams struct {
	ContainerUUID uuid.NullUUID `path:"container_uuid"`
}

func (h *containerHandler) WaitStatus() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *WaitStatusParams) error {
		status := c.Query("status")

		inst, err := h.containerService.Get(c, params.ContainerUUID.UUID)
		if err != nil {
			return err
		}

		return h.containerRunnerService.WaitStatus(c, inst, status)
	})
}
