package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/types"
)

func addSettingsRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetSettings)
	r.PATCH("", handlePatchSettings)
}

func handleGetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, settingsFSAdapter.GetSettings())
}

func handlePatchSettings(c *gin.Context) {
	var settings types.Settings
	err := c.BindJSON(&settings)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	err = settingsService.Update(settings)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}
