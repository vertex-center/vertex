package port

import containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"

type (
	EditorService interface {
		ToYaml(template containerstypes.Template) ([]byte, error)
	}
)
