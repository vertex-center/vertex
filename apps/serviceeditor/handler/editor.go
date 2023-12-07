package handler

import (
	"net/http"

	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/types"
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

func (h *editorHandler) ToYaml(c *router.Context) {
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

func (h *editorHandler) ToYamlInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Convert service to yaml"),
		oapi.Description("Convert service description to a reusable yaml file."),
		oapi.Response(http.StatusOK),
	}
}
