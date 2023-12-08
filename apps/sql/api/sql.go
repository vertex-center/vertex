package api

import (
	"context"

	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/sql/core/types"
)

func (c *Client) GetDBMS(ctx context.Context, containerUuid string) (types.DBMS, error) {
	var dbms types.DBMS
	err := c.Request().
		Pathf("./container/%s", containerUuid).
		ToJSON(&dbms).
		Fetch(ctx)
	return dbms, err
}

func (c *Client) InstallDBMS(ctx context.Context, dbmsId string) (containerstypes.Container, error) {
	var inst containerstypes.Container
	err := c.Request().
		Pathf("./dbms/%s/install", dbmsId).
		Post().
		ToJSON(&inst).
		Fetch(ctx)
	return inst, err
}
