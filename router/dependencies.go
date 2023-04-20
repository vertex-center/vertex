package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/services/dependency"
)

func addDependenciesRoutes(r *gin.RouterGroup) {
	r.POST("/install", handleInstallDependencies)
}

type InstallDependencyBody struct {
	Dependencies []struct {
		Name           string `json:"name"`
		PackageManager string `json:"package_manager"`
	} `json:"dependencies"`
}

func handleInstallDependencies(c *gin.Context) {
	var body InstallDependencyBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	for _, d := range body.Dependencies {
		dep, err := dependency.Get(d.Name)
		if err != nil {
			logger.Warn(fmt.Sprintf("dependency '%s' not found", d.Name))
			continue
		}

		logger.Log(fmt.Sprintf("installing package name='%s' with package_manager=%s", d.Name, d.PackageManager))

		cmd, err := dep.InstallationCommand(d.PackageManager)
		if err != nil {
			logger.Error(err)
			continue
		}

		if cmd.Sudo {
			// Command needs sudo. Sending the command to the client for manual execution.
			c.JSON(http.StatusOK, gin.H{
				"command": cmd.Cmd,
			})
		} else {
			err = cmd.Exec()
			if err != nil {
				logger.Error(err)
				continue
			}
		}

		logger.Log(fmt.Sprintf("package name=%s installed successfully", d.Name))

		c.Status(http.StatusOK)
		return
	}

	_ = c.AbortWithError(http.StatusNotFound, errors.New("failed to find this package manager"))
}
