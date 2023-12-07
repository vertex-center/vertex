package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	DBMSHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		Install() gin.HandlerFunc
		InstallInfo() []oapi.Info
	}
)
