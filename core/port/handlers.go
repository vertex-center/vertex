package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/common/app"
)

type AppsHandler interface {
	GetApps(c *gin.Context) ([]app.Meta, error)
}
