package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/sql"
	sqltypes "github.com/vertex-center/vertex/apps/sql/types"
	"github.com/vertex-center/vertex/types/api"
)

func GetDBMS(ctx context.Context, instanceUuid string) (sqltypes.DBMS, *api.Error) {
	var dbms sqltypes.DBMS
	var apiError api.Error
	err := api.AppRequest(sql.AppRoute).
		Pathf("./instance/%s", instanceUuid).
		ToJSON(&dbms).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return dbms, api.HandleError(err, apiError)
}

func InstallDBMS(ctx context.Context, dbmsId string) (*sqltypes.DBMS, *api.Error) {
	var dbms *sqltypes.DBMS
	var apiError api.Error
	err := api.AppRequest(sql.AppRoute).
		Pathf("./dbms/%s/install", dbmsId).
		Post().
		ToJSON(&dbms).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return dbms, api.HandleError(err, apiError)
}
