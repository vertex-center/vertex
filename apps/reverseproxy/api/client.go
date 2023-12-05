package api

import (
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewReverseProxyClient() *Client {
	port := config.Current.GetPort(reverseproxy.Meta.ID, reverseproxy.Meta.DefaultPort)
	return &Client{
		Client: rest.NewClient(config.Current.Host, port, "/api"),
	}
}
