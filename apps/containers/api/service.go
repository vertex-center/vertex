package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *Client) GetService(ctx context.Context, serviceId string) (types.Service, error) {
	var service types.Service
	err := c.Request().
		Pathf("./service/%s", serviceId).
		ToJSON(&service).
		Fetch(ctx)
	return service, err
}

func (c *Client) InstallService(ctx context.Context, serviceId string) (*types.Container, error) {
	var inst *types.Container
	err := c.Request().
		Pathf("./service/%s/install", serviceId).
		Post().
		ToJSON(&inst).
		Fetch(ctx)
	return inst, err
}
