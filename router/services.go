package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/services/instances"
	servicesmanager "github.com/vertex-center/vertex/services/manager"
)

func addServicesRoutes(r *gin.RouterGroup) {
	r.GET("/available", handleServicesAvailable)
	r.POST("/download", handleServiceDownload)
}

func handleServicesAvailable(c *gin.Context) {
	c.JSON(http.StatusOK, servicesmanager.ListAvailable())
}

type downloadBody struct {
	Repository string `json:"repository"`
}

func handleServiceDownload(c *gin.Context) {
	var body downloadBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	i, err := instances.Install(body.Repository)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"instance": i,
	})
}
