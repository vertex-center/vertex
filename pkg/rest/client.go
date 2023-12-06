package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type Client struct {
	config requests.Config
}

func NewClient(u *url.URL, token string) *Client {
	return &Client{
		config: func(rb *requests.Builder) {
			rb.
				BaseURL(fmt.Sprintf("%s://%s", u.Scheme, u.Host)).
				Path(u.Path + "/").
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
