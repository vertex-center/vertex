package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type tagsHandler struct {
	tagsService port.TagsService
}

func NewTagsHandler(service port.TagsService) port.TagsHandler {
	return &tagsHandler{service}
}

func (h *tagsHandler) GetTags() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (types.Tags, error) {
		return h.tagsService.GetTags(c)
	})
}
