package instancesapi

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

func GetService(ctx context.Context, serviceId string) (types.Service, *types.AppApiError) {
	var service types.Service
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/service/%s", serviceId).
		ToJSON(&service).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return service, types.HandleError(err, apiError)
}

func InstallService(ctx context.Context, serviceId string) (*types.Instance, *types.AppApiError) {
	var inst *types.Instance
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/service/%s/install", serviceId).
		Post().
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return inst, types.HandleError(err, apiError)
}
