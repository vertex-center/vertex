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

func NewClient(host, port, token string) *Client {
	return &Client{
		config: func(rb *requests.Builder) {
			rb.
				BaseURL("http://" + host + ":" + port).
				Path("/api/").
				AddValidator(func(response *http.Response) error {
					log.Request("request from app", vlog.String("path", response.Request.URL.Path))
					return nil
				})

			if token != "" {
				rb.Header("Authorization", "Bearer "+token)
			}
		},
	}
}

func (c *Client) Request() *requests.Builder {
	return requests.New(c.config)
}
