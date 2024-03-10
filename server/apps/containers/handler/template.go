package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/pkg/router"
)

type templateHandler struct {
	containers port.ContainerService
}

func NewTemplateHandler(containerService port.ContainerService) port.TemplateHandler {
	return &templateHandler{containerService}
}

type GetServiceParams struct {
	TemplateID string `path:"template_id"`
}

func (h *templateHandler) GetTemplate() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *GetServiceParams) (*types.Template, error) {
		return h.containers.GetTemplateByID(ctx, params.TemplateID)
	})
}

func (h *templateHandler) GetTemplates() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) ([]types.Template, error) {
		return h.containers.GetTemplates(ctx), nil
	})
}
