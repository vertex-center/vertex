package rest

import (
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type Client struct {
	config requests.Config
}

func NewClient(baseURL string, path string) *Client {
	return &Client{
		config: func(rb *requests.Builder) {
			rb.BaseURL(baseURL).Path(path)
			rb.AddValidator(func(response *http.Response) error {
				log.Request("request from app", vlog.String("path", response.Request.URL.Path))
				return nil
			})
		},
	}
}

func (c *Client) Request() *requests.Builder {
	return requests.New(c.config)
}
