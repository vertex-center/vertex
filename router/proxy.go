package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addProxyRoutes(r *gin.RouterGroup) {
	r.GET("/redirects", handleGetRedirects)
	r.POST("/redirect", handleAddRedirect)
	r.DELETE("/redirect/:id", handleRemoveRedirect)
}

// handleGetRedirects handles the retrieval of all redirects.
func handleGetRedirects(c *gin.Context) {
	redirects := proxyService.GetRedirects()
	c.JSON(http.StatusOK, redirects)
}

type handleAddRedirectBody struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// handleAddRedirect handles the addition of a redirect.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_add_redirect: failed to add the redirect.
func handleAddRedirect(c *gin.Context) {
	var body handleAddRedirectBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	redirect := types.ProxyRedirect{
		Source: body.Source,
		Target: body.Target,
	}

	err = proxyService.AddRedirect(redirect)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToAddRedirect,
			Message: fmt.Sprintf("failed to add redirect: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleRemoveRedirect handles the removal of a redirect.
// Errors can be:
//   - missing_redirect_uuid: missing redirect uuid.
//   - invalid_redirect_uuid: invalid redirect uuid.
//   - failed_to_remove_redirect: failed to remove the redirect.
func handleRemoveRedirect(c *gin.Context) {
	idString := c.Param("id")
	if idString == "" {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrRedirectUuidMissing,
			Message: "missing redirect uuid",
		})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrRedirectUuidInvalid,
			Message: "invalid redirect uuid",
		})
		return
	}

	err = proxyService.RemoveRedirect(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToRemoveRedirect,
			Message: fmt.Sprintf("failed to remove redirect: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
