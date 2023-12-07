package handler

import (
	"errors"
	"net/http"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/core/types/api"
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

func (h *updateHandler) Get(c *router.Context) {
	channel, err := h.settingsService.GetChannel()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetSettings,
			PublicMessage:  "Failed to retrieve update settings.",
			PrivateMessage: err.Error(),
		})
		return
	}

	update, err := h.updateService.GetUpdate(channel)
	if errors.Is(err, types.ErrFailedToFetchBaseline) {
		c.Abort(router.Error{
			Code:           api.ErrFailedToFetchLatestVersion,
			PublicMessage:  "Failed to retrieve latest version information.",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetUpdates,
			PublicMessage:  "Failed to retrieve updates.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(update)
}

func (h *updateHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get the latest update information"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Update{}),
		),
	}
}

func (h *updateHandler) Install(c *router.Context) {
	channel, err := h.settingsService.GetChannel()
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetSettings,
			PublicMessage:  "Failed to retrieve update settings.",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = h.updateService.InstallLatest(channel)
	if errors.Is(err, types.ErrAlreadyUpdating) {
		c.Abort(router.Error{
			Code:           api.ErrAlreadyUpdating,
			PublicMessage:  "Vertex is already Updating. Please wait for the update to finish.",
			PrivateMessage: err.Error(),
		})
		return
	} else if errors.Is(err, types.ErrFailedToFetchBaseline) {
		c.Abort(router.Error{
			Code:           api.ErrFailedToFetchLatestVersion,
			PublicMessage:  "Failed to retrieve latest version information.",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToInstallUpdates,
			PublicMessage:  "Failed to install updates.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *updateHandler) InstallInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Install the latest version"),
		oapi.Description("This endpoint will install the latest version of Vertex."),
		oapi.Response(http.StatusNoContent),
	}
}
