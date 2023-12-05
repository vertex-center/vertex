package api

import (
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewTunnelsClient(token string) *Client {
	port := config.Current.GetPort(tunnels.Meta.ID, tunnels.Meta.DefaultPort)
	return &Client{
		Client: rest.NewClient(config.Current.Host, port, token),
	}
}
