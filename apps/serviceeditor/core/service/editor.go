package service

import (
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"gopkg.in/yaml.v3"
)

type EditorService struct{}

func NewEditorService() port.EditorService {
	return &EditorService{}
}

func (s *EditorService) ToYaml(serv containerstypes.Service) ([]byte, error) {
	return yaml.Marshal(serv)
}
