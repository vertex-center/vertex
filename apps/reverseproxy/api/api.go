package api

import (
	"context"

	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func GetRedirects(ctx context.Context) ([]types.ProxyRedirect, *api.Error) {
	var redirects []types.ProxyRedirect
	var apiError api.Error
	err := api.AppRequest(reverseproxy.AppRoute).
		Path("./redirects").
		ToJSON(&redirects).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return redirects, api.HandleError(err, apiError)
}

func AddRedirect(ctx context.Context, redirect types.ProxyRedirect) *api.Error {
	var apiError api.Error
	err := api.AppRequest(reverseproxy.AppRoute).
		Path("./redirect").
		BodyJSON(&redirect).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func RemoveRedirect(ctx context.Context, id string) *api.Error {
	var apiError api.Error
	err := api.AppRequest(reverseproxy.AppRoute).
		Pathf("./redirect/%s", id).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}
