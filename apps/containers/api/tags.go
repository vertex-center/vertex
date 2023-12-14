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
