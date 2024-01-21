package handler

import (
	"github.com/gin-gonic/gin"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"github.com/vertex-center/vertex/pkg/router"
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
	return router.Handler(func(ctx *gin.Context, template *containerstypes.Template) ([]byte, error) {
		return h.editorService.ToYaml(*template)
	})
}
