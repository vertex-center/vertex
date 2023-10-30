package handler

import (
	"errors"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type UpdateHandler struct {
	updateService   port.UpdateService
	settingsService port.SettingsService
}

func NewUpdateHandler(updateService port.UpdateService, settingsService port.SettingsService) port.UpdateHandler {
	return &UpdateHandler{
		updateService:   updateService,
		settingsService: settingsService,
	}
}

// docapi:begin get_updates
// docapi:method GET
// docapi:summary Get the latest version information.
// docapi:tags updates
// docapi:response 200 Update The latest version information.
// docapi:response 500
// docapi:end

func (h *UpdateHandler) Get(c *router.Context) {
	channel := h.settingsService.GetChannel()

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

// docapi:begin install_update
// docapi:method POST
// docapi:summary Install the latest version.
// docapi:tags updates
// docapi:response 204
// docapi:response 400
// docapi:response 500
// docapi:end

func (h *UpdateHandler) Install(c *router.Context) {
	channel := h.settingsService.GetChannel()

	err := h.updateService.InstallLatest(channel)
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
