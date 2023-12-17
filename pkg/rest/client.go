package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/juju/errors"
)

type Client struct {
	config requests.Config
}

func NewClient(u *url.URL, token, correlationID string) *Client {
	return &Client{
		config: func(rb *requests.Builder) {
			rb.
				BaseURL(fmt.Sprintf("%s://%s", u.Scheme, u.Host)).
				Path(u.Path + "/").
				AddValidator(func(response *http.Response) error {
					if response.StatusCode < 400 {
						return nil
					}

					var buf strings.Builder
					_, err := io.Copy(&buf, response.Body)
					if err != nil {
						return err
					}
					msg := buf.String()

					switch response.StatusCode {
					case http.StatusUnauthorized:
						err = errors.Unauthorized
					case http.StatusForbidden:
						err = errors.Forbidden
					case http.StatusNotFound:
						err = errors.NotFound
					case http.StatusConflict:
						err = errors.AlreadyExists
					case http.StatusUnprocessableEntity:
						err = errors.BadRequest
					default:
						err = errors.New("request error")
					}
					return errors.Annotate(err, msg)
				})

			if token != "" {
				rb.Header("Authorization", "Bearer "+token)
			}
			if correlationID != "" {
				rb.Header("X-Correlation-ID", correlationID)
			}
		},
	}
}

func (c *Client) Request() *requests.Builder {
	return requests.New(c.config)
}
