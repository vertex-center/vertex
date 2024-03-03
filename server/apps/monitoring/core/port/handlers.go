package port

import (
	"github.com/gin-gonic/gin"
)

type (
	MetricsHandler interface {
		GetCollector() gin.HandlerFunc
		InstallCollector() gin.HandlerFunc
		InstallVisualizer() gin.HandlerFunc
	}
)
