package instancesapi

import (
	"context"

	"github.com/carlmjohnson/requests"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

func GetInstances(ctx context.Context) ([]instancestypes.Instance, *types.AppApiError) {
	var instances []instancestypes.Instance
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Path("/api/instances").
		ToJSON(&instances).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return instances, types.HandleError(err, apiError)
}

func CheckForUpdates(ctx context.Context) ([]instancestypes.Instance, *types.AppApiError) {
	var instances []instancestypes.Instance
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Path("/api/instances/checkupdates").
		ToJSON(&instances).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return instances, types.HandleError(err, apiError)
}
