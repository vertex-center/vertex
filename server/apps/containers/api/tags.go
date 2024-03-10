package containersapi

import (
	"context"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
)

func (c *Client) GetTag(ctx context.Context, name string) (types.Tag, error) {
	var tag types.Tag
	err := c.Request().
		Pathf("./tags/%s", name).
		ToJSON(&tag).
		Fetch(ctx)
	return tag, err
}

func (c *Client) GetTags(ctx context.Context) (*types.Tags, error) {
	var tags *types.Tags
	err := c.Request().
		Path("./tags").
		ToJSON(&tags).
		Fetch(ctx)
	return tags, err
}

func (c *Client) CreateTag(ctx context.Context, tag types.Tag) (types.Tag, error) {
	err := c.Request().
		Path("./tags").
		BodyJSON(tag).
		ToJSON(&tag).
		Post().
		Fetch(ctx)
	return tag, err
}

func (c *Client) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return c.Request().
		Pathf("./tag/%s", id).
		Delete().
		Fetch(ctx)
}
