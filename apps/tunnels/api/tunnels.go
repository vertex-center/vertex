package api

import (
	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
)

func InstallTunnel(c *router.Context, provider string) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/tunnels/provider/%s/install", provider).
		Post().
		ErrorJSON(&apiError).
		Fetch(c)
	return types.HandleError(err, apiError)
}
