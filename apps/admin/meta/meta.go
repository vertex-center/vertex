package meta

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/common/app/appmeta"
)

var Meta = appmeta.Meta{
	ID:                "admin",
	Name:              "Vertex Admin",
	Description:       "Administer Vertex",
	Icon:              "admin_panel_settings",
	DefaultPort:       "7500",
	DefaultKernelPort: "7501",
	Dependencies: []*appmeta.Meta{
		&authmeta.Meta,
	},
}
