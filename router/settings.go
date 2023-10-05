package router

import (
	"net/http"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addSettingsRoutes(r *router.Group) {
	r.GET("", handleGetSettings)
	r.PATCH("", handlePatchSettings)
}

// handleGetSettings handles the retrieval of all settings.
func handleGetSettings(c *router.Context) {
	c.JSON(http.StatusOK, settingsFSAdapter.GetSettings())
}

// handlePatchSettings handles the update of all settings.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_update_settings: failed to update the settings.
func handlePatchSettings(c *router.Context) {
	var settings types.Settings
	err := c.ParseBody(&settings)
	if err != nil {
		return
	}

	err = settingsService.Update(settings)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToPatchSettings,
			PublicMessage:  "Failed to update settings.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}
