package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex/services"
)

func InitializeRouter() *gin.Engine {
	r := router.CreateRouter()

	// TODO: Change to POST and read body
	r.GET("/download", handleDownload)
	r.GET("/start", handleStart)
	r.GET("/stop", handleStop)

	return r
}

// Sample service for development purposes
var redisService = services.Service{
	ID:           "vertex-redis",
	Name:         "Vertex Redis",
	Dependencies: []string{},
	Repository:   "github.com/vertex-center/vertex-redis",
}

func handleDownload(c *gin.Context) {
	err := redisService.Download()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleStart(c *gin.Context) {
	runner := services.NewRunner(redisService)

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
