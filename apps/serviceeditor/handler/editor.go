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

// docapi begin vx_devtools_service_editor_to_yaml
// docapi method POST
// docapi summary Convert service to yaml
// docapi desc Convert service description to a reusable yaml file.
// docapi tags Service Editor
// docapi body {Service} The service to convert.
// docapi response 200 {string} The yaml file.
// docapi response 400
// docapi end

func (h *EditorHandler) ToYaml(c *router.Context) {
	var serv containerstypes.Service
	err := c.BindJSON(&serv)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrInvalidService,
			PublicMessage:  "Invalid service.",
			PrivateMessage: "The service is invalid.",
		})
		return
	}

	yaml, err := h.editorService.ToYaml(serv)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrInvalidService,
			PublicMessage:  "Invalid service.",
			PrivateMessage: "The service is invalid.",
		})
		return
	}

	c.Data(200, "yaml", yaml)
}
