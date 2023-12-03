package api

import (
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewAdminClient() *Client {
	return &Client{
		Client: rest.NewClient(config.Current.VertexURL(), "/api/app/admin/"),
	}
}

type KernelClient struct {
	*rest.Client
}

func NewAdminKernelClient() *KernelClient {
	return &KernelClient{
		Client: rest.NewClient(config.Current.KernelURL(), "/api/app/admin/"),
	}
}
