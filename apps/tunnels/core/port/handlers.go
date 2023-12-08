package port

import (
	"github.com/gin-gonic/gin"
)

type ProviderHandler interface {
	Install() gin.HandlerFunc
}
