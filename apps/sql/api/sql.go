package api

import (
	"context"

	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/sql/types"
	"github.com/vertex-center/vertex/types/api"
)

func GetDBMS(ctx context.Context, instanceUuid string) (types.DBMS, *api.Error) {
	var dbms types.DBMS
	var apiError api.Error
	err := api.AppRequest(sql.AppRoute).
		Pathf("./instance/%s", instanceUuid).
		ToJSON(&dbms).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return dbms, api.HandleError(err, apiError)
}

func InstallDBMS(ctx context.Context, dbmsId string) (instancestypes.Instance, *api.Error) {
	var inst instancestypes.Instance
	var apiError api.Error
	err := api.AppRequest(sql.AppRoute).
		Pathf("./dbms/%s/install", dbmsId).
		Post().
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return inst, api.HandleError(err, apiError)
}
