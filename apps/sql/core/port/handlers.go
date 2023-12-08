package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type (
	DBMSHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		Install() gin.HandlerFunc
		InstallInfo() []fizz.OperationOption
	}
)
