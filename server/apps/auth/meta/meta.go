package meta

import (
	logsmeta "github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/common/app/appmeta"
)

var Meta = appmeta.Meta{
	ID:          "auth",
	Name:        "Vertex Auth",
	Description: "Authentication app for Vertex",
	Icon:        "admin_panel_settings",
	Hidden:      true,
	DefaultPort: "7502",
	Dependencies: []*appmeta.Meta{
		&logsmeta.Meta,
	},
}
