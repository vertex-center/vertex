package api

import (
	"context"

	"github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/common/server"
	"github.com/vertex-center/vertex/server/config"
	"github.com/vertex-center/vertex/server/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewLogsClient(ctx context.Context) *Client {
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &Client{
		Client: rest.NewClient(config.Current.Addr(meta.Meta.ID), "", correlationID),
	}
}
