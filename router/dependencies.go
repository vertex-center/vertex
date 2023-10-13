package router

import (
	"fmt"

	types2 "github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addDependenciesRoutes(r *router.Group) {
	r.GET("", handleGetDependencies)
	r.POST("/update", handleUpdateDependencies)
}

// handleGetPackages handles the retrieval of all dependencies.
// Errors can be:
//   - failed_to_get_dependencies: failed to get dependencies.
func handleGetDependencies(c *router.Context) {
	reload := c.Query("reload")

	var dependencies types.Dependencies
	if reload == "true" {
		var err error
		dependencies, err = dependenciesService.CheckForUpdates(settingsService.GetChannel())
		if err != nil {
			c.Abort(router.Error{
				Code:           types2.ErrCodeFailedToCheckForUpdates,
				PublicMessage:  "Failed to check for updates.",
				PrivateMessage: err.Error(),
			})
			return
		}
	} else {
		dependencies = dependenciesService.GetCachedUpdates()
	}

	c.JSON(dependencies)
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
func handleUpdateDependencies(c *router.Context) {
	var body executeUpdatesBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	fmt.Printf("body: %v\n", body)

	var updates []string
	for _, update := range body.Updates {
		updates = append(updates, update.Name)
	}

	err = dependenciesService.InstallUpdates(updates)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToInstallUpdates,
			PublicMessage:  "Failed to install updates.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
