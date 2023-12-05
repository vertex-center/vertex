package meta

import "github.com/vertex-center/vertex/core/types/app"

var Meta = app.Meta{
	ID:                "containers",
	Name:              "Vertex Containers",
	Description:       "Create and manage containers.",
	Icon:              "deployed_code",
	DefaultPort:       "7504",
	DefaultKernelPort: "7505",
}
