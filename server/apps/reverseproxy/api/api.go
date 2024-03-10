package api

import (
	"context"

	"github.com/vertex-center/vertex/server/apps/reverseproxy/core/types"
)

func (c *Client) GetRedirects(ctx context.Context) ([]types.ProxyRedirect, error) {
	var redirects []types.ProxyRedirect
	err := c.Request().
		Path("./redirects").
		ToJSON(&redirects).
		Fetch(ctx)
	return redirects, err
}

func (c *Client) AddRedirect(ctx context.Context, redirect types.ProxyRedirect) error {
	return c.Request().
		Path("./redirect").
		BodyJSON(&redirect).
		Post().
		Fetch(ctx)
}

func (c *Client) RemoveRedirect(ctx context.Context, id string) error {
	return c.Request().
		Pathf("./redirect/%s", id).
		Delete().
		Fetch(ctx)
}
