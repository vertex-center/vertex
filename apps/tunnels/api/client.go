package api

import (
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewTunnelsClient(token string) *Client {
	return &Client{
		Client: rest.NewClient(tunnels.Meta.ApiURL(), token),
	}
}
