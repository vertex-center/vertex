package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

func addPackagesRoutes(r *gin.RouterGroup) {
	r.POST("/install", handleInstallPackages)
}

type InstallPackagesBody struct {
	Packages []struct {
		Name           string `json:"name"`
		PackageManager string `json:"package_manager"`
	} `json:"packages"`
}

// handleInstallPackages handles the installation of packages.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
func handleInstallPackages(c *gin.Context) {
	var body InstallPackagesBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	for _, d := range body.Packages {
		pkg, err := packageService.GetByID(d.Name)
		if err != nil {
			log.Warn("dependency not found",
				vlog.String("name", d.Name),
				vlog.String("package_manager", d.PackageManager),
			)
			continue
		}

		log.Info("installing package",
			vlog.String("name", d.Name),
			vlog.String("package_manager", d.PackageManager),
		)

		cmd, err := packageService.InstallationCommand(&pkg, d.PackageManager)
		if err != nil {
			log.Error(err)
			continue
		}

		if cmd.Sudo {
			// Command needs sudo. Sending the command to the client for manual execution.
			c.JSON(http.StatusOK, gin.H{
				"command": cmd.Cmd,
			})
		} else {
			err = packageService.Install(cmd)
			if err != nil {
				log.Error(err)
				continue
			}
		}

		log.Info("package installed successfully",
			vlog.String("name", d.Name),
			vlog.String("package_manager", d.PackageManager),
		)

		c.Status(http.StatusNoContent)
		return
	}

	c.Status(http.StatusNoContent)
}
