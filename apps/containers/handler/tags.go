package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/types/session"
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

type GetTagParams struct {
	Name string `query:"name"`
}

func (h *tagsHandler) GetTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetTagParams) (*types.Tag, error) {
		s := session.Get(c)
		tag, err := h.tagsService.GetTag(c, s.UserID, params.Name)
		if err != nil {
			return nil, err
		}
		return &tag, nil
	})
}

func (h *tagsHandler) GetTags() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (types.Tags, error) {
		s := session.Get(c)
		return h.tagsService.GetTags(c, s.UserID)
	})
}

type CreateTagParams struct {
	Tag types.Tag `json:"tag"`
}

func (h *tagsHandler) CreateTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *CreateTagParams) (types.Tag, error) {
		s := session.Get(c)
		params.Tag.UserID = s.UserID
		return h.tagsService.CreateTag(c, params.Tag)
	})
}

type DeleteTagParams struct {
	ID types.TagID `path:"id"`
}

func (h *tagsHandler) DeleteTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteTagParams) error {
		return h.tagsService.DeleteTag(c, params.ID)
	})
}
