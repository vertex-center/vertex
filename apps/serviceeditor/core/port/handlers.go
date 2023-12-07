package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	EditorHandler interface {
		ToYaml() gin.HandlerFunc
		ToYamlInfo() []oapi.Info
	}
)
