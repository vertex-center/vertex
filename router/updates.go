package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addUpdatesRoutes(r *gin.RouterGroup, currentVertexVersion string) {
	r.GET("", func(c *gin.Context) {
		handleGetUpdates(c, currentVertexVersion)
	})

	r.POST("", func(c *gin.Context) {
		handleExecuteUpdates(c, currentVertexVersion)
	})
}

func handleGetUpdates(c *gin.Context, currentVertexVersion string) {
	updates, err := updateService.CheckForUpdates(currentVertexVersion)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updates)
}

type executeUpdatesBody struct {
	Updates []struct {
		Name string `json:"name"`
	}
}

func handleExecuteUpdates(c *gin.Context, currentVertexVersion string) {
	var body executeUpdatesBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	var errors []error

	for _, update := range body.Updates {
		var err error

		switch update.Name {
		case "vertex":
			err = updateService.InstallVertexUpdate(currentVertexVersion)
		default:
			logger.Error(fmt.Errorf("the service name %s is not valid for an update", update.Name))
		}

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		for _, err := range errors {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.Status(http.StatusOK)
}
