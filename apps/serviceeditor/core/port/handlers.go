package port

import (
	"github.com/gin-gonic/gin"
)

type (
	EditorHandler interface {
		ToYaml() gin.HandlerFunc
	}
)
