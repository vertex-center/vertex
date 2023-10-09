package instancesapi

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

func GetInstances(ctx context.Context) ([]types.Instance, *types.AppApiError) {
	var instances []types.Instance
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Path("/api/instances").
		ToJSON(&instances).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return instances, types.HandleError(err, apiError)
}

func CheckForUpdates(ctx context.Context) ([]types.Instance, *types.AppApiError) {
	var instances []types.Instance
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Path("/api/instances/checkupdates").
		ToJSON(&instances).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return instances, types.HandleError(err, apiError)
}
