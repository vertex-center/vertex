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

func handleGetUpdates(c *gin.Context) {
	reload := c.Query("reload")

	var updates types.Updates
	if reload == "true" {
		var err error
		updates, err = updateService.CheckForUpdates()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
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

func handleExecuteUpdates(c *gin.Context) {
	var body executeUpdatesBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	var updates []string
	for _, update := range body.Updates {
		updates = append(updates, update.Name)
	}

	err = updateService.InstallUpdates(updates)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Status(http.StatusOK)
}
