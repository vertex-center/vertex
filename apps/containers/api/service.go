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
