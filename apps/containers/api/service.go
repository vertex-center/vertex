package containersapi

import (
	"context"
	"github.com/vertex-center/vertex/apps/containers"
	types2 "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func GetService(ctx context.Context, serviceId string) (types2.Service, *api.Error) {
	var service types2.Service
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./service/%s", serviceId).
		ToJSON(&service).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return service, api.HandleError(err, apiError)
}

func InstallService(ctx context.Context, serviceId string) (*types2.Container, *api.Error) {
	var inst *types2.Container
	var apiError api.Error
	err := api.AppRequest(containers.AppRoute).
		Pathf("./service/%s/install", serviceId).
		Post().
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return inst, api.HandleError(err, apiError)
}
