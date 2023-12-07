package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type updateHandler struct {
	updateService   port.UpdateService
	settingsService port.SettingsService
}

func NewUpdateHandler(updateService port.UpdateService, settingsService port.SettingsService) port.UpdateHandler {
	return &updateHandler{
		updateService:   updateService,
		settingsService: settingsService,
	}
}

func (h *updateHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*types.Update, error) {
		channel, err := h.settingsService.GetChannel()
		if err != nil {
			return nil, err
		}
		update, err := h.updateService.GetUpdate(channel)
		if err != nil {
			return nil, err
		}
		return update, nil
	})
}

func (h *updateHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getUpdate"),
		oapi.Summary("Get the latest update information"),
	}
}

func (h *updateHandler) Install() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) error {
		channel, err := h.settingsService.GetChannel()
		if err != nil {
			return err
		}
		return h.updateService.InstallLatest(channel)
	})
}

func (h *updateHandler) InstallInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("installUpdate"),
		oapi.Summary("Install the latest version"),
		oapi.Description("This endpoint will install the latest version of Vertex."),
	}
}
