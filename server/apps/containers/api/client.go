package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/common/server"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewContainersClient(ctx context.Context) *Client {
	token := ctx.Value("token").(string)
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &Client{
		Client: rest.NewClient(config.Current.Addr(meta.Meta.ID), token, correlationID),
	}
}

type KernelClient struct {
	*rest.Client
}

func NewContainersKernelClient(ctx context.Context) *KernelClient {
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &KernelClient{
		Client: rest.NewClient(config.Current.KernelAddr(meta.Meta.ID), "", correlationID),
	}
}
