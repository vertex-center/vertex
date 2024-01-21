package service

import (
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"gopkg.in/yaml.v3"
)

type editorService struct{}

func NewEditorService() port.EditorService {
	return &editorService{}
}

func (s *editorService) ToYaml(template containerstypes.Template) ([]byte, error) {
	return yaml.Marshal(template)
}
