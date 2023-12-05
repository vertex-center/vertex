package api

import (
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewMonitoringClient(token string) *Client {
	return &Client{
		Client: rest.NewClient(monitoring.Meta.ApiURL(), token),
	}
}
