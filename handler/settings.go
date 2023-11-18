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

// docapi begin get_settings
// docapi method GET
// docapi summary Get settings.
// docapi tags Settings
// docapi response 200 {Settings} The settings.
// docapi end

func (h *SettingsHandler) Get(c *router.Context) {
	c.JSON(h.settingsService.Get())
}

// docapi begin patch_settings
// docapi method PATCH
// docapi summary Update settings.
// docapi tags Settings
// docapi body {Settings} The settings to update.
// docapi response 200
// docapi response 400
// docapi response 500
// docapi end

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
