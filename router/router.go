package router

import (
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

	// TODO: Change to POST and read body
	r.GET("/download", handleDownload)
	r.GET("/start", handleStart)
	r.GET("/stop", handleStop)

	r.GET("/installed", handleInstalled)
	r.GET("/available", handleAvailable)

	return r
}

func handleDownload(c *gin.Context) {
	err := servicesmanager.Download(servicesmanager.ListAvailable()[0])
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleStart(c *gin.Context) {
	runner := services.NewRunner(servicesmanager.ListAvailable()[0])

	err := runner.Start()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleStop(c *gin.Context) {
	err := services.GetRunner().Stop()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
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
