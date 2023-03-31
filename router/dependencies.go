package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/dependency"
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
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	for _, d := range body.Dependencies {
		dep, err := dependency.Get(d.Name)
		if err != nil {
			logger.Warn(fmt.Sprintf("dependency '%s' not found", d.Name))
			continue
		}

		logger.Log(fmt.Sprintf("installing package name='%s' with package_manager=%s", d.Name, d.PackageManager))

		err = dep.Install(d.PackageManager)
		if err != nil {
			logger.Error(err)
			continue
		}

		logger.Log(fmt.Sprintf("package name=%s installed successfully", d.Name))
	}

	c.Status(http.StatusOK)
}
