package meta

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/common/app"
)

var Meta = app.Meta{
	ID:                "containers",
	Name:              "Vertex Containers",
	Description:       "Create and manage containers.",
	Icon:              "deployed_code",
	DefaultPort:       "7504",
	DefaultKernelPort: "7505",
	Dependencies: []*app.Meta{
		&authmeta.Meta,
	},
}
