package api

import (
	"github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewAuthClient(token string) *Client {
	port := config.Current.GetPort(meta.Meta.ID, meta.Meta.DefaultPort)
	return &Client{
		Client: rest.NewClient(config.Current.Host, port, token),
	}
}
