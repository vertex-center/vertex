package instancesapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/instances"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/types/api"
)

func GetInstances(ctx context.Context) ([]types.Instance, *api.Error) {
	var insts []types.Instance
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Path("./instances").
		ToJSON(&insts).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return insts, api.HandleError(err, apiError)
}

func CheckForUpdates(ctx context.Context) ([]types.Instance, *api.Error) {
	var insts []types.Instance
	var apiError api.Error
	err := api.AppRequest(instances.AppRoute).
		Path("./instances/checkupdates").
		ToJSON(&insts).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return insts, api.HandleError(err, apiError)
}
