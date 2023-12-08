package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
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

func (h *updateHandler) GetInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("getUpdate"),
		fizz.Summary("Get the latest update information"),
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

func (h *updateHandler) InstallInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("installUpdate"),
		fizz.Summary("Install the latest version"),
		fizz.Description("This endpoint will install the latest version of Vertex."),
	}
}
