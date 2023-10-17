package handler

import (
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type EditorHandler struct {
	editorService port.EditorService
}

func NewEditorHandler(editorService port.EditorService) port.EditorHandler {
	return &EditorHandler{
		editorService: editorService,
	}
}

func (h *EditorHandler) ToYaml(c *router.Context) {
	var serv containerstypes.Service
	err := c.BindJSON(&serv)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrInvalidService,
			PublicMessage:  "Invalid service.",
			PrivateMessage: "The service is invalid.",
		})
	}

	yaml, err := h.editorService.ToYaml(serv)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrInvalidService,
			PublicMessage:  "Invalid service.",
			PrivateMessage: "The service is invalid.",
		})
	}

	c.Data(200, "yaml", yaml)
}
