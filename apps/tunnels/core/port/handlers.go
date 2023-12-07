package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type ProviderHandler interface {
	Install(c *router.Context)
	InstallInfo() []oapi.Info
}
