package handler

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/core/types/api"
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

// docapi begin get_settings
// docapi method GET
// docapi summary Get settings
// docapi tags Settings
// docapi response 200 {Settings} The settings.
// docapi response 500
// docapi end

func (h *settingsHandler) Get(c *router.Context) {
	settings, err := h.service.Get()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetSettings,
			PublicMessage:  "Failed to get settings.",
			PrivateMessage: err.Error(),
		})
		return
	}
	c.JSON(settings)
}

// docapi begin patch_settings
// docapi method PATCH
// docapi summary Patch settings
// docapi tags Settings
// docapi body {Settings} The settings to update.
// docapi response 200
// docapi response 400
// docapi response 500
// docapi end

func (h *settingsHandler) Patch(c *router.Context) {
	var settings types.AdminSettings
	err := c.ParseBody(&settings)
	if err != nil {
		return
	}

	err = h.service.Update(settings)
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
