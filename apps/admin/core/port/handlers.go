package port

import (
	"github.com/gin-gonic/gin"
)

type (
	ChecksHandler interface {
		Check() gin.HandlerFunc
	}

	DatabaseHandler interface {
		GetCurrentDbms() gin.HandlerFunc
		MigrateTo() gin.HandlerFunc
	}

	SettingsHandler interface {
		Get() gin.HandlerFunc
		Patch() gin.HandlerFunc
	}

	UpdateHandler interface {
		Get() gin.HandlerFunc
	}
)
