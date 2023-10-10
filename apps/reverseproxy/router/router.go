package router

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/reverseproxy/adapter"
	"github.com/vertex-center/vertex/apps/reverseproxy/service"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

type AppRouter struct {
	proxyFSAdapter types.ProxyAdapterPort

	proxyService *service.ProxyService

	proxyRouter *ProxyRouter
}

func NewAppRouter() *AppRouter {
	r := &AppRouter{
		proxyFSAdapter: adapter.NewProxyFSAdapter(nil),
	}
	r.proxyService = service.NewProxyService(r.proxyFSAdapter)

	return r
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	group.GET("/redirects", r.handleGetRedirects)
	group.POST("/redirect", r.handleAddRedirect)
	group.DELETE("/redirect/:id", r.handleRemoveRedirect)
}

func (r *AppRouter) GetServices() []types.AppService {
	return []types.AppService{
		r.proxyService,
	}
}

func (r *AppRouter) GetProxyService() *service.ProxyService {
	return r.proxyService
}

// handleGetRedirects handles the retrieval of all redirects.
func (r *AppRouter) handleGetRedirects(c *router.Context) {
	redirects := r.proxyService.GetRedirects()
	c.JSON(redirects)
}

type handleAddRedirectBody struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// handleAddRedirect handles the addition of a redirect.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_add_redirect: failed to add the redirect.
func (r *AppRouter) handleAddRedirect(c *router.Context) {
	var body handleAddRedirectBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	redirect := types.ProxyRedirect{
		Source: body.Source,
		Target: body.Target,
	}

	err = r.proxyService.AddRedirect(redirect)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToAddRedirect,
			PublicMessage:  fmt.Sprintf("Failed to add redirect '%s' to '%s'.", redirect.Source, redirect.Target),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// handleRemoveRedirect handles the removal of a redirect.
// Errors can be:
//   - missing_redirect_uuid: missing redirect uuid.
//   - invalid_redirect_uuid: invalid redirect uuid.
//   - failed_to_remove_redirect: failed to remove the redirect.
func (r *AppRouter) handleRemoveRedirect(c *router.Context) {
	idString := c.Param("id")
	if idString == "" {
		c.BadRequest(router.Error{
			Code:           api.ErrRedirectUuidMissing,
			PublicMessage:  "The request is missing the redirect UUID.",
			PrivateMessage: "Field 'id' is required.",
		})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           api.ErrRedirectUuidInvalid,
			PublicMessage:  "The redirect UUID is invalid.",
			PrivateMessage: err.Error(),
		})
		return
	}

	err = r.proxyService.RemoveRedirect(id)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToRemoveRedirect,
			PublicMessage:  fmt.Sprintf("Failed to remove redirect '%s'.", id),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
