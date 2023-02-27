package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex/services"
	servicesmanager "github.com/vertex-center/vertex/services/manager"
	"github.com/vertex-center/vertex/services/runners"
)

func InitializeRouter() *gin.Engine {
	r := router.CreateRouter()
	r.Use(cors.Default())

	servicesGroup := r.Group("/services")
	servicesGroup.GET("/installed", handleInstalled)
	servicesGroup.GET("/available", handleAvailable)
	servicesGroup.POST("/download", handleDownload)

	serviceGroup := r.Group("/service/:service_id")
	serviceGroup.POST("/start", handleStart)
	serviceGroup.POST("/stop", handleStop)

	return r
}

func handleInstalled(c *gin.Context) {
	installed, err := servicesmanager.ListInstalled()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, installed)
}

func handleAvailable(c *gin.Context) {
	c.JSON(http.StatusOK, servicesmanager.ListAvailable())
}

type DownloadBody struct {
	Service services.Service `json:"service"`
}

func handleDownload(c *gin.Context) {
	var body DownloadBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	err = servicesmanager.Download(body.Service)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleStart(c *gin.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("service_id was missing in the URL"))
		return
	}

	runner, err := runners.NewRunner(serviceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = runner.Start()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleStop(c *gin.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("service_id was missing in the URL"))
		return
	}

	runner, err := runners.GetRunner(serviceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = runner.Stop()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
