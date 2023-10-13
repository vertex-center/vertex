package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/types/api"
)

func GetService(ctx context.Context, serviceId string) (types.Service, *api.Error) {
	var service types.Service
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./service/%s", serviceId).
		ToJSON(&service).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return service, api.HandleError(err, apiError)
}

func InstallService(ctx context.Context, serviceId string) (*types.Container, *api.Error) {
	var inst *types.Container
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./service/%s/install", serviceId).
		Post().
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return inst, api.HandleError(err, apiError)
}
