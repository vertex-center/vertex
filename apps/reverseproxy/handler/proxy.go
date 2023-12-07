package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type proxyHandler struct {
	proxyService port.ProxyService
}

func NewProxyHandler(proxyService port.ProxyService) port.ProxyHandler {
	return &proxyHandler{
		proxyService: proxyService,
	}
}

func (r *proxyHandler) GetRedirects() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (types.ProxyRedirects, error) {
		return r.proxyService.GetRedirects(), nil
	})
}

func (r *proxyHandler) GetRedirectsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getRedirects"),
		oapi.Summary("Get redirects"),
	}
}

func (r *proxyHandler) AddRedirect() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, redirect *types.ProxyRedirect) error {
		err := r.proxyService.AddRedirect(*redirect)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to add redirect '%s' to '%s'", redirect.Source, redirect.Target))
		}
		return nil
	})
}

func (r *proxyHandler) AddRedirectInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("addRedirect"),
		oapi.Summary("Add redirect"),
	}
}

type RemoveRedirectParams struct {
	ID uuid.UUID `path:"id"`
}

func (r *proxyHandler) RemoveRedirect() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *RemoveRedirectParams) error {
		err := r.proxyService.RemoveRedirect(params.ID)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to remove redirect '%s'", params.ID))
		}
		return nil
	})
}

func (r *proxyHandler) RemoveRedirectInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("removeRedirect"),
		oapi.Summary("Remove redirect"),
	}
}
