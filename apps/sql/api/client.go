package api

import (
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewSqlClient() *Client {
	port := config.Current.GetPort(sql.Meta.ID, sql.Meta.DefaultPort)
	return &Client{
		Client: rest.NewClient(config.Current.Host, port, "/api"),
	}
}
