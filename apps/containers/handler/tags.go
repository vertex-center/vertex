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
	return router.Handler(func(ctx *gin.Context, params *GetTagParams) (*types.Tag, error) {
		s := session.Get(ctx)
		tag, err := h.tagsService.GetTag(ctx, s.UserID, params.Name)
		if err != nil {
			return nil, err
		}
		return &tag, nil
	})
}

func (h *tagsHandler) GetTags() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) (types.Tags, error) {
		s := session.Get(ctx)
		return h.tagsService.GetTags(ctx, s.UserID)
	})
}

type CreateTagParams struct {
	Name string `json:"name"`
}

func (h *tagsHandler) CreateTag() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *CreateTagParams) (types.Tag, error) {
		s := session.Get(ctx)
		return h.tagsService.CreateTag(ctx, types.Tag{
			UserID: s.UserID,
			Name:   params.Name,
		})
	})
}

type DeleteTagParams struct {
	ID uuid.NullUUID `path:"id"`
}

func (h *tagsHandler) DeleteTag() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *DeleteTagParams) error {
		return h.tagsService.DeleteTag(ctx, params.ID.UUID)
	})
}
