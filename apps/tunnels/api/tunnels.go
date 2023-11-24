package api

import (
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

func (c *Client) InstallTunnel(ctx *router.Context, provider string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./provider/%s/install", provider).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}
