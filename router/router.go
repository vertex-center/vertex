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
)

func InitializeRouter() *gin.Engine {
	r := router.CreateRouter()
	r.Use(cors.Default())

	servicesGroup := r.Group("/services")
	servicesGroup.GET("", handleServicesInstalled)
	servicesGroup.GET("/available", handleServicesAvailable)
	servicesGroup.POST("/download", handleServiceDownload)

	serviceGroup := r.Group("/service/:service_id")
	serviceGroup.POST("/start", handleServiceStart)
	serviceGroup.POST("/stop", handleServiceStop)

	return r
}

func handleServicesInstalled(c *gin.Context) {
	installed := services.ListInstalled()
	c.JSON(http.StatusOK, installed)
}

func handleServicesAvailable(c *gin.Context) {
	c.JSON(http.StatusOK, servicesmanager.ListAvailable())
}

type DownloadBody struct {
	Service services.Service `json:"service"`
}

func handleServiceDownload(c *gin.Context) {
	var body DownloadBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	service, err := body.Service.Install()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"service": service,
	})
}

func handleServiceStart(c *gin.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("service_id was missing in the URL"))
		return
	}

	service, err := services.GetInstalled(serviceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = service.Start()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleServiceStop(c *gin.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("service_id was missing in the URL"))
		return
	}

	service, err := services.GetInstalled(serviceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = service.Stop()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
