package api

import (
	"github.com/vertex-center/vertex/apps/admin/meta"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewAdminClient(token string) *Client {
	return &Client{
		Client: rest.NewClient(meta.Meta.ApiURL(), token),
	}
}

type KernelClient struct {
	*rest.Client
}

func NewAdminKernelClient() *KernelClient {
	return &KernelClient{
		Client: rest.NewClient(meta.Meta.ApiKernelURL(), ""),
	}
}
