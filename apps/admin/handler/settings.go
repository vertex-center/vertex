package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type settingsHandler struct {
	service port.SettingsService
}

func NewSettingsHandler(settingsService port.SettingsService) port.SettingsHandler {
	return &settingsHandler{
		service: settingsService,
	}
}

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

func (h *settingsHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get settings"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.AdminSettings{}),
		),
	}
}

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

func (h *settingsHandler) PatchInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Patch settings"),
		oapi.Response(http.StatusNoContent),
	}
}
