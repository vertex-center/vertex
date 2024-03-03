package api

import (
	"github.com/gin-gonic/gin"
)

func (c *Client) InstallTunnel(ctx *gin.Context, provider string) error {
	return c.Request().
		Pathf("./provider/%s/install", provider).
		Post().
		Fetch(ctx)
}
