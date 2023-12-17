package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/uuid"
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
	Name string `json:"name"`
}

func (h *tagsHandler) CreateTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *CreateTagParams) (types.Tag, error) {
		s := session.Get(c)
		return h.tagsService.CreateTag(c, types.Tag{
			UserID: s.UserID,
			Name:   params.Name,
		})
	})
}

type DeleteTagParams struct {
	ID uuid.NullUUID `path:"id"`
}

func (h *tagsHandler) DeleteTag() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteTagParams) error {
		return h.tagsService.DeleteTag(c, params.ID.UUID)
	})
}
