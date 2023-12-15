package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *Client) GetTags(ctx context.Context) (*types.Tags, error) {
	var tags *types.Tags
	err := c.Request().
		Path("./tags").
		ToJSON(&tags).
		Fetch(ctx)
	return tags, err
}

func (c *Client) CreateTag(ctx context.Context, tag types.Tag) error {
	return c.Request().
		Path("./tags").
		BodyJSON(tag).
		Post().
		Fetch(ctx)
}

func (c *Client) DeleteTag(ctx context.Context, id types.TagID) error {
	return c.Request().
		Pathf("./tags/%s", id).
		Delete().
		Fetch(ctx)
}
