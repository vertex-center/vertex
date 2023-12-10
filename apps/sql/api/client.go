package api

import (
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewSqlClient(token string) *Client {
	return &Client{
		Client: rest.NewClient(config.Current.URL(sql.Meta.ID), token),
	}
}
