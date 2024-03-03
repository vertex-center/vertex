package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
	"github.com/vertex-center/vertex/pkg/router"
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
	return router.Handler(func(ctx *gin.Context) (types.ProxyRedirects, error) {
		return r.proxyService.GetRedirects(), nil
	})
}

func (r *proxyHandler) AddRedirect() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, redirect *types.ProxyRedirect) error {
		err := r.proxyService.AddRedirect(*redirect)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to add redirect '%s' to '%s'", redirect.Source, redirect.Target))
		}
		return nil
	})
}

type RemoveRedirectParams struct {
	ID uuid.NullUUID `path:"id"`
}

func (r *proxyHandler) RemoveRedirect() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *RemoveRedirectParams) error {
		err := r.proxyService.RemoveRedirect(params.ID.UUID)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to remove redirect '%s'", params.ID.UUID))
		}
		return nil
	})
}
