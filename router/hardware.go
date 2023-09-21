package router

import (
	"github.com/gin-gonic/gin"
)

func addHardwareRoutes(r *gin.RouterGroup) {
	r.GET("", handleGetHardware)
}

func handleGetHardware(c *gin.Context) {
	c.JSON(200, hardwareService.GetHardware())
}
