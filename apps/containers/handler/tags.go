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

type GetTagParams struct {
	Name string `query:"name"`
}

func (h *tagsHandler) GetTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetTagParams) (*types.Tag, error) {
		userID := uint(c.GetInt("user_id"))
		tag, err := h.tagsService.GetTag(c, userID, params.Name)
		if err != nil {
			return nil, err
		}
		return &tag, nil
	})
}

func (h *tagsHandler) GetTags() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (types.Tags, error) {
		userID := uint(c.GetInt("user_id"))
		return h.tagsService.GetTags(c, userID)
	})
}

type CreateTagParams struct {
	Tag types.Tag `json:"tag"`
}

func (h *tagsHandler) CreateTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *CreateTagParams) (types.Tag, error) {
		userID := uint(c.GetInt("user_id"))
		params.Tag.UserID = userID
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
