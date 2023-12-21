package meta

import (
	logsmeta "github.com/vertex-center/vertex/apps/logs/meta"
	"github.com/vertex-center/vertex/common/app/appmeta"
)

var Meta = appmeta.Meta{
	ID:          "client",
	Name:        "Vertex Client",
	Description: "Vertex web client",
	Icon:        "web",
	DefaultPort: "7518",
	Hidden:      true,
	Dependencies: []*appmeta.Meta{
		&logsmeta.Meta,
	},
}
