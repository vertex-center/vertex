package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type (
	MetricsHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		InstallCollector() gin.HandlerFunc
		InstallCollectorInfo() []fizz.OperationOption

		InstallVisualizer() gin.HandlerFunc
		InstallVisualizerInfo() []fizz.OperationOption
	}
)
