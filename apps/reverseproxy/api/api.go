package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) GetRedirects(ctx context.Context) ([]types.ProxyRedirect, *api.Error) {
	var redirects []types.ProxyRedirect
	var apiError api.Error
	err := c.Request().
		Path("./redirects").
		ToJSON(&redirects).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return redirects, api.HandleError(err, apiError)
}

func (c *Client) AddRedirect(ctx context.Context, redirect types.ProxyRedirect) *api.Error {
	var apiError api.Error
	err := c.Request().
		Path("./redirect").
		BodyJSON(&redirect).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) RemoveRedirect(ctx context.Context, id string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./redirect/%s", id).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}
