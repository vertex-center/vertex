package meta

import "github.com/vertex-center/vertex/core/types/app"

var Meta = app.Meta{
	ID:                "admin",
	Name:              "Vertex Admin",
	Description:       "Administer Vertex",
	Icon:              "admin_panel_settings",
	DefaultPort:       "7500",
	DefaultKernelPort: "7501",
}
