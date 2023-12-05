package rest

import (
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type Client struct {
	config requests.Config
}

func NewClient(host, port, path string) *Client {
	return &Client{
		config: func(rb *requests.Builder) {
			rb.BaseURL(host + ":" + port).Path(path)
			rb.AddValidator(func(response *http.Response) error {
				log.Request("request from app", vlog.String("path", response.Request.URL.Path))
				return nil
			})
			rb.Header("Authorization", "Bearer "+config.Current.MasterApiKey)
		},
	}
}

func (c *Client) Request() *requests.Builder {
	return requests.New(c.config)
}
