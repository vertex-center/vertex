package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *KernelClient) GetContainers(ctx context.Context) ([]types.DockerContainer, *api.Error) {
	var apiError api.Error
	var containers []types.DockerContainer
	err := c.Request().
		Pathf("./docker/containers").
		ErrorJSON(&apiError).
		ToJSON(&containers).
		Fetch(ctx)
	return containers, api.HandleError(err, apiError)
}
