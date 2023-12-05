package containersapi

import (
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewContainersClient() *Client {
	port := config.Current.GetPort(containers.Meta.ID, containers.Meta.DefaultPort)
	return &Client{
		Client: rest.NewClient(config.Current.Host, port, "/api"),
	}
}
