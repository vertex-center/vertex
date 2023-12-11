package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/common/server"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewSqlClient(ctx context.Context) *Client {
	token := ctx.Value("token").(string)
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &Client{
		Client: rest.NewClient(config.Current.URL(sql.Meta.ID), token, correlationID),
	}
}
