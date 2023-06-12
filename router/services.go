package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addServicesRoutes(r *gin.RouterGroup) {
	r.GET("/", handleGetService)
	r.GET("/available", handleServicesAvailable)
	r.POST("/download", handleServiceDownload)
}

func handleGetService(c *gin.Context) {
	repo := c.Query("repository")

	service, err := serviceService.Get(repo)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, service)
}

func handleServicesAvailable(c *gin.Context) {
	c.JSON(http.StatusOK, serviceService.ListAvailable())
}

type downloadBody struct {
	Repository  string `json:"repository"`
	UseDocker   *bool  `json:"use_docker,omitempty"`
	UseReleases *bool  `json:"use_releases,omitempty"`
}

func handleServiceDownload(c *gin.Context) {
	var body downloadBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	i, err := instanceService.Install(body.Repository, body.UseDocker, body.UseReleases)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"instance": i,
	})
}
