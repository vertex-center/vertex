package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/common/app"
	"github.com/wI2L/fizz"
)

type (
	AppsHandler interface {
		GetApps(c *gin.Context) ([]app.Meta, error)
		GetAppsInfo() []fizz.OperationOption
	}

	DebugHandler interface {
		HardReset(c *gin.Context) error
		HardResetInfo() []fizz.OperationOption
	}
)
