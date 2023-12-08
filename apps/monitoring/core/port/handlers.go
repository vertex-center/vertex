package port

import (
	"github.com/gin-gonic/gin"
)

type (
	MetricsHandler interface {
		Get() gin.HandlerFunc
		InstallCollector() gin.HandlerFunc
		InstallVisualizer() gin.HandlerFunc
	}
)
