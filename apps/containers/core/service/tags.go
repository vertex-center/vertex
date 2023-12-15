package service

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type tagsService struct {
	tags port.TagAdapter
}

func NewTagsService(tags port.TagAdapter) port.TagsService {
	return &tagsService{tags}
}

func (s *tagsService) GetTags(ctx context.Context) (types.Tags, error) {
	return s.tags.GetTags(ctx)
}

func (s *tagsService) CreateTag(ctx context.Context, tag types.Tag) error {
	return s.tags.CreateTag(ctx, tag)
}

func (s *tagsService) DeleteTag(ctx context.Context, id types.TagID) error {
	return s.tags.DeleteTag(ctx, id)
}
