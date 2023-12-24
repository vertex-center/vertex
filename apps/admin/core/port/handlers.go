package port

import (
	"github.com/gin-gonic/gin"
)

type (
	ChecksHandler interface {
		Check() gin.HandlerFunc
	}

	SettingsHandler interface {
		Get() gin.HandlerFunc
		Patch() gin.HandlerFunc
	}

	UpdateHandler interface {
		Get() gin.HandlerFunc
	}
)
