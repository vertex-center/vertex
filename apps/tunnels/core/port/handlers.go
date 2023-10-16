package port

import "github.com/vertex-center/vertex/pkg/router"

type ProviderHandler interface {
	Install(c *router.Context)
}
