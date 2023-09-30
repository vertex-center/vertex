package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addSettingsRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetSettings)
	r.PATCH("", handlePatchSettings)
}

// handleGetSettings handles the retrieval of all settings.
func handleGetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, settingsFSAdapter.GetSettings())
}

// handlePatchSettings handles the update of all settings.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_update_settings: failed to update the settings.
func handlePatchSettings(c *gin.Context) {
	var settings types.Settings
	err := c.BindJSON(&settings)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	err = settingsService.Update(settings)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToPatchSettings,
			Message: fmt.Sprintf("failed to update settings: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}
