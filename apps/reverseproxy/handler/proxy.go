package handler

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type ProxyHandler struct {
	proxyService port.ProxyService
}

func NewProxyHandler(proxyService port.ProxyService) port.ProxyHandler {
	return &ProxyHandler{
		proxyService: proxyService,
	}
}

// docapi begin vx_reverse_proxy_get_redirects
// docapi method GET
// docapi summary Get redirects
// docapi tags Apps/Reverse Proxy
// docapi response 200 {[]Redirect} The redirects.
// docapi end

func (r *ProxyHandler) GetRedirects(c *router.Context) {
	redirects := r.proxyService.GetRedirects()
	c.JSON(redirects)
}

type AddRedirectBody struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// docapi begin vx_reverse_proxy_add_redirect
// docapi method POST
// docapi summary Add redirect
// docapi tags Apps/Reverse Proxy
// docapi body {AddRedirectBody} The redirect to add.
// docapi response 204
// docapi response 400
// docapi response 500
// docapi end

func (r *ProxyHandler) AddRedirect(c *router.Context) {
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

// docapi begin vx_reverse_proxy_remove_redirect
// docapi method DELETE
// docapi summary Remove redirect
// docapi tags Apps/Reverse Proxy
// docapi query id {string} The redirect UUID.
// docapi response 204
// docapi response 400
// docapi response 500
// docapi end

func (r *ProxyHandler) RemoveRedirect(c *router.Context) {
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
