package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	AppsHandler interface {
		GetApps(c *gin.Context) ([]app.Meta, error)
		GetAppsInfo() []oapi.Info
	}

	DebugHandler interface {
		HardReset(c *gin.Context) error
		HardResetInfo() []oapi.Info
	}
)
