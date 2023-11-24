package containersapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) GetContainers(ctx context.Context) (map[uuid.UUID]*types.Container, *api.Error) {
	var insts map[uuid.UUID]*types.Container
	var apiError api.Error
	err := c.Request().
		Path("./containers").
		ToJSON(&insts).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return insts, api.HandleError(err, apiError)
}

func (c *Client) CheckForUpdates(ctx context.Context) ([]types.Container, *api.Error) {
	var insts []types.Container
	var apiError api.Error
	err := c.Request().
		Path("./containers/checkupdates").
		ToJSON(&insts).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return insts, api.HandleError(err, apiError)
}
