package api

import (
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
)

func InstallTunnel(c *router.Context, provider string) *api.Error {
	var apiError api.Error
	err := api.AppRequest(tunnels.AppRoute).
		Pathf("./provider/%s/install", provider).
		Post().
		ErrorJSON(&apiError).
		Fetch(c)
	return api.HandleError(err, apiError)
}
