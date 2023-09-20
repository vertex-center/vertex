package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/types"
)

func addDependenciesRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetDependencies)
	r.POST("/update", handleUpdateDependencies)
}

// handleGetPackages handles the retrieval of all dependencies.
// Errors can be:
//   - failed_to_get_dependencies: failed to get dependencies.
func handleGetDependencies(c *gin.Context) {
	reload := c.Query("reload")

	var dependencies types.Dependencies
	if reload == "true" {
		var err error
		dependencies, err = dependenciesService.CheckForUpdates()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIError{
				Code:    "failed_to_check_for_updates",
				Message: fmt.Sprintf("failed to check for updates: %v", err),
			})
			return
		}
	} else {
		dependencies = dependenciesService.GetCachedUpdates()
	}

	c.JSON(http.StatusOK, dependencies)
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
func handleUpdateDependencies(c *gin.Context) {
	var body executeUpdatesBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	fmt.Printf("body: %v\n", body)

	var updates []string
	for _, update := range body.Updates {
		updates = append(updates, update.Name)
	}

	err = dependenciesService.InstallUpdates(updates)
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
