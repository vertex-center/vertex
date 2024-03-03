package meta

import (
	"github.com/vertex-center/vertex/common/app/appmeta"
)

var Meta = appmeta.Meta{
	ID:          "logs",
	Name:        "Vertex Logs",
	Description: "Gather and view logs from all Vertex apps",
	Icon:        "feed",
	DefaultPort: "7516",
}
