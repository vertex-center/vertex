package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type (
	EditorHandler interface {
		ToYaml() gin.HandlerFunc
		ToYamlInfo() []fizz.OperationOption
	}
)
