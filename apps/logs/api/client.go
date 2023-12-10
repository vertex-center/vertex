package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/logs/meta"
	"github.com/vertex-center/vertex/common/server"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewLogsClient(ctx context.Context) *Client {
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &Client{
		Client: rest.NewClient(config.Current.URL(meta.Meta.ID), "", correlationID),
	}
}
