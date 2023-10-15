package router

import (
	"errors"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addUpdateRoutes(r *router.Group) {
	r.GET("", handleGetLatestUpdate)
	r.POST("", handleInstallLatestUpdate)
}

func handleGetLatestUpdate(c *router.Context) {
	channel := settingsService.GetChannel()

	update, err := updateService.GetUpdate(channel)
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

func handleInstallLatestUpdate(c *router.Context) {
	channel := settingsService.GetChannel()

	err := updateService.InstallLatest(channel)
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
