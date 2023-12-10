package api

import (
	"github.com/vertex-center/vertex/apps/admin/meta"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewAdminClient(token string) *Client {
	return &Client{
		Client: rest.NewClient(config.Current.URL(meta.Meta.ID), token),
	}
}

type KernelClient struct {
	*rest.Client
}

func NewAdminKernelClient() *KernelClient {
	return &KernelClient{
		Client: rest.NewClient(config.Current.KernelURL(meta.Meta.ID), ""),
	}
}
