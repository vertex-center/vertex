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
	return &Client{
		Client: rest.NewClient(config.Current.URL(meta.Meta.ID), token),
	}
}
