package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type (
	ProxyHandler interface {
		GetRedirects() gin.HandlerFunc
		GetRedirectsInfo() []fizz.OperationOption

		AddRedirect() gin.HandlerFunc
		AddRedirectInfo() []fizz.OperationOption

		RemoveRedirect() gin.HandlerFunc
		RemoveRedirectInfo() []fizz.OperationOption
	}
)
