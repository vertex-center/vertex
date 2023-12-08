package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/core/types"
)

func (c *KernelClient) GetContainers(ctx context.Context) ([]types.DockerContainer, error) {
	var containers []types.DockerContainer
	err := c.Request().
		Pathf("./docker/containers").
		ToJSON(&containers).
		Fetch(ctx)
	return containers, err
}
