package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type ProviderHandler interface {
	Install() gin.HandlerFunc
	InstallInfo() []oapi.Info
}
