package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/server/apps/admin/core/port"
	"github.com/vertex-center/vertex/server/apps/admin/core/types"
	"github.com/vertex-center/vertex/server/pkg/router"
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
	return router.Handler(func(ctx *gin.Context) (*types.Update, error) {
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
