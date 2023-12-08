package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type settingsHandler struct {
	service port.SettingsService
}

func NewSettingsHandler(settingsService port.SettingsService) port.SettingsHandler {
	return &settingsHandler{
		service: settingsService,
	}
}

func (h *settingsHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*types.AdminSettings, error) {
		settings, err := h.service.Get()
		if err != nil {
			return nil, err
		}
		return &settings, nil
	})
}

type PatchSettingsParams struct {
	Settings types.AdminSettings `json:"settings"`
}

func (h *settingsHandler) Patch() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchSettingsParams) error {
		return h.service.Update(params.Settings)
	})
}
