package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addServicesRoutes(r *gin.RouterGroup) {
	r.GET("/available", handleServicesAvailable)
	r.POST("/install", handleServiceInstall)
}

func handleServicesAvailable(c *gin.Context) {
	c.JSON(http.StatusOK, serviceService.ListAvailable())
}

type downloadBody struct {
	Method    string `json:"method"`
	ServiceID string `json:"service_id"`
}

func handleServiceInstall(c *gin.Context) {
	var body downloadBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	i, err := instanceService.Install(body.ServiceID, body.Method)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"instance": i,
	})
}
