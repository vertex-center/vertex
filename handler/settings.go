package handler

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type SettingsHandler struct {
	settingsService port.SettingsService
}

func NewSettingsHandler(settingsService port.SettingsService) port.SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
	}
}

func (h *SettingsHandler) Get(c *router.Context) {
	c.JSON(h.settingsService.Get())
}

func (h *SettingsHandler) Patch(c *router.Context) {
	var settings types.Settings
	err := c.ParseBody(&settings)
	if err != nil {
		return
	}

	err = h.settingsService.Update(settings)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToPatchSettings,
			PublicMessage:  "Failed to update settings.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
