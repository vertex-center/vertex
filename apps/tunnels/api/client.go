package api

import (
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewTunnelsClient() *Client {
	return &Client{
		Client: rest.NewClient(config.Current.VertexURL(), "/api/app/tunnels/"),
	}
}
