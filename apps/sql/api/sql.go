package api

import (
	"context"

	"github.com/carlmjohnson/requests"
	sqltypes "github.com/vertex-center/vertex/apps/sql/types"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

func GetDBMS(ctx context.Context, instanceUuid string) (sqltypes.DBMS, *types.AppApiError) {
	var dbms sqltypes.DBMS
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/sql/%s", instanceUuid).
		ToJSON(&dbms).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return dbms, types.HandleError(err, apiError)
}

func InstallDBMS(ctx context.Context, dbmsId string) (*sqltypes.DBMS, *types.AppApiError) {
	var dbms *sqltypes.DBMS
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/sql/db/%s/install", dbmsId).
		Post().
		ToJSON(&dbms).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return dbms, types.HandleError(err, apiError)
}
