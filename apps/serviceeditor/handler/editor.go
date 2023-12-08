package handler

import (
	"github.com/gin-gonic/gin"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
)

type editorHandler struct {
	editorService port.EditorService
}

func NewEditorHandler(editorService port.EditorService) port.EditorHandler {
	return &editorHandler{
		editorService: editorService,
	}
}

func (h *editorHandler) ToYaml() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, serv *containerstypes.Service) ([]byte, error) {
		return h.editorService.ToYaml(*serv)
	})
}

func (h *editorHandler) ToYamlInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("toYaml"),
		fizz.Summary("Convert service to yaml"),
		fizz.Description("Convert service description to a reusable yaml file."),
	}
}
