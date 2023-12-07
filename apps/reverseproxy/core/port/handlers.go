package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	ProxyHandler interface {
		GetRedirects() gin.HandlerFunc
		GetRedirectsInfo() []oapi.Info

		AddRedirect() gin.HandlerFunc
		AddRedirectInfo() []oapi.Info

		RemoveRedirect() gin.HandlerFunc
		RemoveRedirectInfo() []oapi.Info
	}
)
