package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) InstallTunnel(ctx *gin.Context, provider string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./provider/%s/install", provider).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}
