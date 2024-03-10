package service

import (
	"context"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
)

type tagsService struct {
	tags port.TagAdapter
}

func NewTagsService(tags port.TagAdapter) port.TagsService {
	return &tagsService{tags}
}

func (s *tagsService) GetTag(ctx context.Context, userID uuid.UUID, name string) (types.Tag, error) {
	return s.tags.GetTag(ctx, userID, name)
}

func (s *tagsService) GetTags(ctx context.Context, userID uuid.UUID) (types.Tags, error) {
	return s.tags.GetTags(ctx, userID)
}

func (s *tagsService) CreateTag(ctx context.Context, tag types.Tag) (types.Tag, error) {
	tag.ID = uuid.New()
	err := s.tags.CreateTag(ctx, tag)
	return tag, err
}

func (s *tagsService) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return s.tags.DeleteTag(ctx, id)
}
