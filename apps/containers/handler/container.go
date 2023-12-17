package handler

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

type containerHandler struct {
	ctx              *apptypes.Context
	containerService port.ContainerService
}

func NewContainerHandler(ctx *apptypes.Context, containerService port.ContainerService) port.ContainerHandler {
	return &containerHandler{
		ctx:              ctx,
		containerService: containerService,
	}
}

type GetContainerParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) Get() gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, params *GetContainerParams) (*types.Container, error) {
		return h.containerService.Get(c, params.ContainerID)
	}, http.StatusOK)
}

type DeleteContainerParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) Delete() gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, params *DeleteContainerParams) error {
		return h.containerService.Delete(c, params.ContainerID)
	}, http.StatusNoContent)
}

type PatchBodyDatabase struct {
	ContainerID  types.ContainerID `json:"container_id"`
	DatabaseName *string           `json:"db_name"`
}

type PatchContainerParams struct {
	ContainerID     types.ContainerID `path:"container_id"`
	LaunchOnStartup *bool             `json:"launch_on_startup,omitempty"`
	Name            *string           `json:"name,omitempty"`
	ImageTag        *string           `json:"image_tag,omitempty"`
}

func (h *containerHandler) Patch() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerID)
		if err != nil {
			return err
		}

		if params.LaunchOnStartup != nil {
			inst.LaunchOnStartup = *params.LaunchOnStartup
		}
		if params.Name != nil && *params.Name != "" {
			inst.Name = *params.Name
		}
		if params.ImageTag != nil {
			inst.ImageTag = *params.ImageTag
		}

		return h.containerService.UpdateContainer(c, params.ContainerID, *inst)
	})
}

type StartContainerParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) Start() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StartContainerParams) error {
		return h.containerService.Start(c, params.ContainerID)
	})
}

type StopContainerParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) Stop() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *StopContainerParams) error {
		return h.containerService.Stop(c, params.ContainerID)
	})
}

type AddTagParams struct {
	ContainerID types.ContainerID `path:"container_id"`
	TagID       types.TagID       `path:"tag_id"`
}

func (h *containerHandler) AddContainerTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *AddTagParams) error {
		return h.containerService.AddContainerTag(c, params.ContainerID, params.TagID)
	})
}

type GetContainerEnvParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) GetContainerEnv() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetContainerEnvParams) (types.EnvVariables, error) {
		return h.containerService.GetContainerEnv(c, params.ContainerID)
	})
}

type PatchEnvironmentParams struct {
	ContainerID types.ContainerID  `path:"container_id"`
	Env         types.EnvVariables `body:"env"`
}

func (h *containerHandler) PatchEnvironment() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchEnvironmentParams) error {
		return h.containerService.SaveEnv(c, params.ContainerID, params.Env)
	})
}

type EventsContainerParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) Events() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *EventsContainerParams) error {
		inst, err := h.containerService.Get(c, params.ContainerID)
		if err != nil {
			return err
		}

		eventsChan := make(chan sse.Event)
		defer close(eventsChan)

		done := c.Request.Context().Done()

		listener := event.NewTempListener(func(e event.Event) error {
			switch e := e.(type) {
			case types.EventContainerLog:
				if inst.ID != e.ContainerID {
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
				if inst.ID != e.ContainerID {
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
		return h.containerService.GetContainerInfo(c, params.ContainerID)
	})
}

type RecreateContainerParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) RecreateDocker() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *RecreateContainerParams) error {
		return h.containerService.RecreateContainer(c, params.ContainerID)
	})
}

type LogsContainerParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) GetLogs() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *LogsContainerParams) ([]types.LogLine, error) {
		return h.containerService.GetLatestLogs(params.ContainerID)
	})
}

type GetVersionsParams struct {
	ContainerID types.ContainerID `path:"container_id"`
	UseCache    bool              `query:"cache"`
}

func (h *containerHandler) GetVersions() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetVersionsParams) ([]string, error) {
		log.Info("GetVersions", vlog.Bool("use_cache", params.UseCache))
		return h.containerService.GetAllVersions(c, params.ContainerID, params.UseCache)
	})
}

type WaitStatusParams struct {
	ContainerID types.ContainerID `path:"container_id"`
}

func (h *containerHandler) WaitStatus() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *WaitStatusParams) error {
		status := c.Query("status")
		return h.containerService.WaitStatus(c, params.ContainerID, status)
	})
}
