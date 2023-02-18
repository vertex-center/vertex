package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/quentinguidee/installer-service/services"
	"github.com/quentinguidee/microservice-core/router"
	"net/http"
)

func InitializeRouter() *gin.Engine {
	r := router.CreateRouter()

	// TODO: Change to POST and read body
	r.GET("/download", handleDownload)

	return r
}

func handleDownload(c *gin.Context) {
	// Sample service for development purposes
	service := services.Service{
		ID:           "redis-service",
		Name:         "Redis Service",
		Dependencies: []string{},
		Repository:   "github.com/quentinguidee/redis-service",
	}

	err := service.Download()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
