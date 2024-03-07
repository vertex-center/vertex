package handler

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/event"
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
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) Get() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *GetContainerParams) (*types.Container, error) {
		return h.containerService.Get(ctx, params.ContainerID.UUID)
	}, http.StatusOK)
}

func (h *containerHandler) GetContainers() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context) (types.Containers, error) {
		filters := types.ContainerFilters{}

		features := ctx.QueryArray("features[]")
		if len(features) > 0 {
			filters.Features = &features
		}

		tags := ctx.QueryArray("tags[]")
		if len(tags) > 0 {
			filters.Tags = &tags
		}

		return h.containerService.GetContainersWithFilters(ctx, filters)
	}, http.StatusOK)
}

type CreateContainerParams struct {
	TemplateID *string `json:"template_id,omitempty"`
	Image      *string `json:"image,omitempty"`
	ImageTag   *string `json:"image_tag,omitempty"`
}

func (h *containerHandler) CreateContainer() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *CreateContainerParams) (*types.Container, error) {
		return h.containerService.CreateContainer(ctx, types.CreateContainerOptions{
			TemplateID: params.TemplateID,
			Image:      params.Image,
			ImageTag:   params.ImageTag,
		})
	}, http.StatusCreated)
}

type DeleteContainerParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) Delete() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *DeleteContainerParams) error {
		return h.containerService.Delete(ctx, params.ContainerID.UUID)
	}, http.StatusOK)
}

type PatchBodyDatabase struct {
	ContainerID  uuid.NullUUID `json:"container_id"`
	DatabaseName *string       `json:"db_name"`
}

type PatchContainerParams struct {
	ContainerID     uuid.NullUUID `path:"container_id"`
	LaunchOnStartup *bool         `json:"launch_on_startup,omitempty"`
	Name            *string       `json:"name,omitempty"`
	ImageTag        *string       `json:"image_tag,omitempty"`
}

func (h *containerHandler) Patch() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *PatchContainerParams) error {
		c, err := h.containerService.Get(ctx, params.ContainerID.UUID)
		if err != nil {
			return err
		}

		if params.LaunchOnStartup != nil {
			c.LaunchOnStartup = *params.LaunchOnStartup
		}
		if params.Name != nil && *params.Name != "" {
			c.Name = *params.Name
		}
		if params.ImageTag != nil {
			c.ImageTag = *params.ImageTag
		}

		return h.containerService.UpdateContainer(ctx, params.ContainerID.UUID, *c)
	}, http.StatusOK)
}

type StartContainerParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) Start() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *StartContainerParams) error {
		return h.containerService.Start(ctx, params.ContainerID.UUID)
	}, http.StatusOK)
}

type StopContainerParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) Stop() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *StopContainerParams) error {
		return h.containerService.Stop(ctx, params.ContainerID.UUID)
	}, http.StatusOK)
}

type AddTagParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
	TagID       uuid.NullUUID `path:"tag_id"`
}

func (h *containerHandler) AddContainerTag() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *AddTagParams) error {
		return h.containerService.AddContainerTag(ctx, params.ContainerID.UUID, params.TagID.UUID)
	}, http.StatusOK)
}

func (h *containerHandler) GetDocker() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *GetContainerParams) (map[string]any, error) {
		return h.containerService.GetContainerInfo(ctx, params.ContainerID.UUID)
	}, http.StatusOK)
}

type RecreateContainerParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) RecreateDocker() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *RecreateContainerParams) error {
		return h.containerService.RecreateContainer(ctx, params.ContainerID.UUID)
	}, http.StatusNoContent)
}

type LogsContainerParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) GetLogs() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *LogsContainerParams) ([]types.LogLine, error) {
		return h.containerService.GetLatestLogs(params.ContainerID.UUID)
	}, http.StatusOK)
}

type GetVersionsParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
	UseCache    bool          `query:"cache"`
}

func (h *containerHandler) GetVersions() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *GetVersionsParams) ([]string, error) {
		log.Info("GetVersions", vlog.Bool("use_cache", params.UseCache))
		return h.containerService.GetAllVersions(ctx, params.ContainerID.UUID, params.UseCache)
	}, http.StatusOK)
}

type WaitStatusParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) WaitStatus() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *WaitStatusParams) error {
		status := ctx.Query("status")
		return h.containerService.WaitStatus(ctx, params.ContainerID.UUID, status)
	}, http.StatusNoContent)
}

func (h *containerHandler) CheckForUpdates() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context) (types.Containers, error) {
		return h.containerService.CheckForUpdates(ctx)
	}, http.StatusOK)
}

type EventsContainerParams struct {
	ContainerID uuid.NullUUID `path:"container_id"`
}

func (h *containerHandler) ContainerEvents() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context, params *EventsContainerParams) error {
		inst, err := h.containerService.Get(ctx, params.ContainerID.UUID)
		if err != nil {
			return err
		}

		eventsChan := make(chan sse.Event)
		defer close(eventsChan)

		done := ctx.Request.Context().Done()

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

		ctx.Stream(func(w io.Writer) bool {
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
	}, http.StatusOK)
}

func (h *containerHandler) ContainersEvents() gin.HandlerFunc {
	return tonic.Handler(func(ctx *gin.Context) error {
		eventsChan := make(chan sse.Event)
		defer close(eventsChan)

		done := ctx.Request.Context().Done()

		listener := event.NewTempListener(func(e event.Event) error {
			switch e.(type) {
			case types.EventContainersChange:
				eventsChan <- sse.Event{
					Event: types.EventNameContainersChange,
				}
			}
			return nil
		})

		h.ctx.AddListener(listener)
		defer h.ctx.RemoveListener(listener)

		first := true

		ctx.Stream(func(w io.Writer) bool {
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
	}, http.StatusOK)
}
