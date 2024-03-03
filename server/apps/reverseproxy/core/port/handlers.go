package port

import (
	"github.com/gin-gonic/gin"
)

type (
	ProxyHandler interface {
		GetRedirects() gin.HandlerFunc
		AddRedirect() gin.HandlerFunc
		RemoveRedirect() gin.HandlerFunc
	}
)
