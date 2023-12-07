package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func (r *proxyHandler) GetRedirects(c *router.Context) {
	redirects := r.proxyService.GetRedirects()
	c.JSON(redirects)
}

func (r *proxyHandler) GetRedirectsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get redirects"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.ProxyRedirects{}),
		),
	}
}

type AddRedirectBody struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (r *proxyHandler) AddRedirect(c *router.Context) {
	var body AddRedirectBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	redirect := types.ProxyRedirect{
		Source: body.Source,
		Target: body.Target,
	}

	err = r.proxyService.AddRedirect(redirect)
	if errors.Is(err, types.ErrSourceInvalid) {
		c.Abort(router.Error{
			Code:           types.ErrCodeSourceInvalid,
			PublicMessage:  "Failed to add redirect as source is empty.",
			PrivateMessage: err.Error(),
		})
		return
	} else if errors.Is(err, types.ErrTargetInvalid) {
		c.Abort(router.Error{
			Code:           types.ErrCodeTargetInvalid,
			PublicMessage:  "Failed to add redirect as target is empty.",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToAddRedirect,
			PublicMessage:  fmt.Sprintf("Failed to add redirect '%s' to '%s'.", redirect.Source, redirect.Target),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (r *proxyHandler) AddRedirectInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Add redirect"),
		oapi.Response(http.StatusNoContent),
	}
}

func (r *proxyHandler) RemoveRedirect(c *router.Context) {
	idString := c.Param("id")
	if idString == "" {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeRedirectUuidMissing,
			PublicMessage:  "The request is missing the redirect UUID.",
			PrivateMessage: "Field 'id' is required.",
		})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeRedirectUuidInvalid,
			PublicMessage:  "The redirect UUID is invalid.",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = r.proxyService.RemoveRedirect(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToRemoveRedirect,
			PublicMessage:  fmt.Sprintf("Failed to remove redirect '%s'.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (r *proxyHandler) RemoveRedirectInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Remove redirect"),
		oapi.Response(http.StatusNoContent),
	}
}
