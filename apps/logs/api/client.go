package api

import (
	"github.com/vertex-center/vertex/apps/logs/meta"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewLogsClient() *Client {
	return &Client{
		Client: rest.NewClient(config.Current.URL(meta.Meta.ID), ""),
	}
}
