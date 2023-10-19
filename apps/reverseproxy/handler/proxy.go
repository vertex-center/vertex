package handler

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	types2 "github.com/vertex-center/vertex/apps/reverseproxy/core/types"

	"github.com/google/uuid"
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

func (r *ProxyHandler) GetRedirects(c *router.Context) {
	redirects := r.proxyService.GetRedirects()
	c.JSON(redirects)
}

type AddRedirectBody struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (r *ProxyHandler) AddRedirect(c *router.Context) {
	var body AddRedirectBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	redirect := types2.ProxyRedirect{
		Source: body.Source,
		Target: body.Target,
	}

	err = r.proxyService.AddRedirect(redirect)
	if errors.Is(err, types2.ErrSourceInvalid) {
		c.Abort(router.Error{
			Code:           types2.ErrCodeSourceInvalid,
			PublicMessage:  fmt.Sprint("Failed to add redirect as source is empty."),
			PrivateMessage: err.Error(),
		})
		return
	} else if errors.Is(err, types2.ErrTargetInvalid) {
		c.Abort(router.Error{
			Code:           types2.ErrCodeTargetInvalid,
			PublicMessage:  fmt.Sprint("Failed to add redirect as target is empty."),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types2.ErrCodeFailedToAddRedirect,
			PublicMessage:  fmt.Sprintf("Failed to add redirect '%s' to '%s'.", redirect.Source, redirect.Target),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (r *ProxyHandler) RemoveRedirect(c *router.Context) {
	idString := c.Param("id")
	if idString == "" {
		c.BadRequest(router.Error{
			Code:           types2.ErrCodeRedirectUuidMissing,
			PublicMessage:  "The request is missing the redirect UUID.",
			PrivateMessage: "Field 'id' is required.",
		})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types2.ErrCodeRedirectUuidInvalid,
			PublicMessage:  "The redirect UUID is invalid.",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = r.proxyService.RemoveRedirect(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           types2.ErrCodeFailedToRemoveRedirect,
			PublicMessage:  fmt.Sprintf("Failed to remove redirect '%s'.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
