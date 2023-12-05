package api

import (
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewReverseProxyClient(token string) *Client {
	return &Client{
		Client: rest.NewClient(reverseproxy.Meta.ApiURL(), token),
	}
}
