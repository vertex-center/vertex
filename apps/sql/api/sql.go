package api

import (
	"context"

	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/sql/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) GetDBMS(ctx context.Context, containerUuid string) (types.DBMS, *api.Error) {
	var dbms types.DBMS
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s", containerUuid).
		ToJSON(&dbms).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return dbms, api.HandleError(err, apiError)
}

func (c *Client) InstallDBMS(ctx context.Context, dbmsId string) (containerstypes.Container, *api.Error) {
	var inst containerstypes.Container
	var apiError api.Error
	err := c.Request().
		Pathf("./dbms/%s/install", dbmsId).
		Post().
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return inst, api.HandleError(err, apiError)
}
