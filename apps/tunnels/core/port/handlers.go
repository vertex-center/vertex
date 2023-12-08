package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type ProviderHandler interface {
	Install() gin.HandlerFunc
	InstallInfo() []fizz.OperationOption
}
