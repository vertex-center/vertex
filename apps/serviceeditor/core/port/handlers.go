package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	EditorHandler interface {
		ToYaml(c *router.Context)
		ToYamlInfo() []oapi.Info
	}
)
