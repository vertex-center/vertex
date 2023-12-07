package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	MetricsHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		InstallCollector() gin.HandlerFunc
		InstallCollectorInfo() []oapi.Info

		InstallVisualizer() gin.HandlerFunc
		InstallVisualizerInfo() []oapi.Info
	}
)
