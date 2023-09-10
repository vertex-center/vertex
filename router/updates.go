package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/types"
)

func addUpdatesRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetUpdates)
	r.POST("", handleExecuteUpdates)
}

// handleGetUpdates handles the retrieval of all updates.
// Errors can be:
//   - failed_to_check_for_updates: failed to check for updates.
func handleGetUpdates(c *gin.Context) {
	reload := c.Query("reload")

	var updates types.Updates
	if reload == "true" {
		var err error
		updates, err = updateService.CheckForUpdates()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIError{
				Code:    "failed_to_check_for_updates",
				Message: fmt.Sprintf("failed to check for updates: %v", err),
			})
			return
		}
	} else {
		updates = updateService.GetCachedUpdates()
	}

	c.JSON(http.StatusOK, updates)
}

type executeUpdatesBody struct {
	Updates []struct {
		Name string `json:"name"`
	}
}

// handleExecuteUpdates handles the execution of updates.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_install_updates: failed to install the updates.
//   - failed_to_reload_services: failed to reload the services.
//   - failed_to_reload_packages: failed to reload the packages.
func handleExecuteUpdates(c *gin.Context) {
	var body executeUpdatesBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	var updates []string
	for _, update := range body.Updates {
		updates = append(updates, update.Name)
	}

	err = updateService.InstallUpdates(updates)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_install_updates",
			Message: fmt.Sprintf("failed to install updates: %v", err),
		})
		return
	}

	err = serviceService.Reload()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_reload_services",
			Message: fmt.Sprintf("failed to reload services: %v", err),
		})
		return
	}

	err = packageService.Reload()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_reload_packages",
			Message: fmt.Sprintf("failed to reload packages: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
