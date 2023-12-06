package containersapi

import (
	"github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewContainersClient(token string) *Client {
	return &Client{
		Client: rest.NewClient(meta.Meta.ApiURL(), token),
	}
}

type KernelClient struct {
	*rest.Client
}

func NewContainersKernelClient() *KernelClient {
	return &KernelClient{
		Client: rest.NewClient(meta.Meta.ApiKernelURL(), ""),
	}
}
