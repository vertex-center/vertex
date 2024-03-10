package api

import (
	"context"

	"github.com/vertex-center/vertex/server/apps/tunnels"
	"github.com/vertex-center/vertex/server/common/server"
	"github.com/vertex-center/vertex/server/config"
	"github.com/vertex-center/vertex/server/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewTunnelsClient(ctx context.Context) *Client {
	token := ctx.Value("token").(string)
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &Client{
		Client: rest.NewClient(config.Current.Addr(tunnels.Meta.ID), token, correlationID),
	}
}
