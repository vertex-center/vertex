package port

import (
	"github.com/gin-gonic/gin"
)

type (
	DBMSHandler interface {
		Get() gin.HandlerFunc
		Install() gin.HandlerFunc
	}
)
