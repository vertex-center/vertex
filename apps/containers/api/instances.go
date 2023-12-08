package containersapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *Client) GetContainers(ctx context.Context) (map[uuid.UUID]*types.Container, error) {
	var insts map[uuid.UUID]*types.Container
	err := c.Request().
		Path("./containers").
		ToJSON(&insts).
		Fetch(ctx)
	return insts, err
}

func (c *Client) CheckForUpdates(ctx context.Context) ([]types.Container, error) {
	var insts []types.Container
	err := c.Request().
		Path("./containers/checkupdates").
		ToJSON(&insts).
		Fetch(ctx)
	return insts, err
}
