package port

import containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"

type (
	EditorService interface {
		ToYaml(serv containerstypes.Service) ([]byte, error)
	}
)
