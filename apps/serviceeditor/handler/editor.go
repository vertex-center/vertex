package handler

import (
	"github.com/gin-gonic/gin"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
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

func (h *editorHandler) ToYamlInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("toYaml"),
		oapi.Summary("Convert service to yaml"),
		oapi.Description("Convert service description to a reusable yaml file."),
	}
}
