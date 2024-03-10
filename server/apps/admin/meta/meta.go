package meta

import (
	authmeta "github.com/vertex-center/vertex/server/apps/auth/meta"
	logsmeta "github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/common/app/appmeta"
)

var Meta = appmeta.Meta{
	ID:          "admin",
	Name:        "Vertex Admin",
	Description: "Administer Vertex",
	Icon:        "admin_panel_settings",
	DefaultPort: "7500",
	Dependencies: []*appmeta.Meta{
		&authmeta.Meta,
		&logsmeta.Meta,
	},
}
