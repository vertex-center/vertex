package api

import (
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewMonitoringClient(token string) *Client {
	port := config.Current.GetPort(monitoring.Meta.ID, monitoring.Meta.DefaultPort)
	return &Client{
		Client: rest.NewClient(config.Current.Host, port, token),
	}
}
