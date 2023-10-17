package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	EditorHandler interface {
		ToYaml(c *router.Context)
	}
)
