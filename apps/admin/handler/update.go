package handler

import (
	"errors"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
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

// docapi begin get_updates
// docapi method GET
// docapi summary Get the latest version info
// docapi tags Updates
// docapi response 200 {Update} The latest version information.
// docapi response 500
// docapi end

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

// docapi begin install_update
// docapi method POST
// docapi summary Install the latest version
// docapi tags Updates
// docapi response 204
// docapi response 400
// docapi response 500
// docapi end

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
