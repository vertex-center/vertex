package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/logger"
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

func handleInstallPackages(c *gin.Context) {
	var body InstallPackagesBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	for _, d := range body.Packages {
		pkg, err := packageService.Get(d.Name)
		if err != nil {
			logger.Warn("dependency not found").
				AddKeyValue("name", d.Name).
				AddKeyValue("package_manager", d.PackageManager).
				Print()

			continue
		}

		logger.Log(fmt.Sprintf("installing package name='%s' with package_manager=%s", d.Name, d.PackageManager)).Print()

		cmd, err := packageService.InstallationCommand(&pkg, d.PackageManager)
		if err != nil {
			logger.Error(err).Print()
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
				logger.Error(err).Print()
				continue
			}
		}

		logger.Log("package installed successfully").
			AddKeyValue("name", d.Name).
			AddKeyValue("package_manager", d.PackageManager).
			Print()

		c.Status(http.StatusOK)
		return
	}

	_ = c.AbortWithError(http.StatusNotFound, errors.New("failed to find this package manager"))
}
