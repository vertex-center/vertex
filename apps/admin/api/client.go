package api

import (
	"github.com/vertex-center/vertex/apps/admin/meta"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewAdminClient() *Client {
	port := config.Current.GetPort(meta.Meta.ID, meta.Meta.DefaultPort)
	return &Client{
		Client: rest.NewClient(config.Current.Host, port, "/api"),
	}
}

type KernelClient struct {
	*rest.Client
}

func NewAdminKernelClient() *KernelClient {
	port := config.Current.GetPort(meta.Meta.ID, meta.Meta.DefaultKernelPort)
	return &KernelClient{
		Client: rest.NewClient(config.Current.Host, port, "/api"),
	}
}
