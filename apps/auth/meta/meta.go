package meta

import (
	"github.com/vertex-center/vertex/common/app"
)

var Meta = app.Meta{
	ID:          "auth",
	Name:        "Vertex Auth",
	Description: "Authentication app for Vertex",
	Icon:        "admin_panel_settings",
	Hidden:      true,
	DefaultPort: "7502",
}
