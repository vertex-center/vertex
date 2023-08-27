package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

func addProxyRoutes(r *gin.RouterGroup) {
	r.GET("/redirects", handleGetRedirects)
	r.POST("/redirect", handleAddRedirect)
	r.DELETE("/redirect/:id", handleRemoveRedirect)
}

func handleGetRedirects(c *gin.Context) {
	redirects := proxyService.GetRedirects()
	c.JSON(http.StatusOK, redirects)
}

type handleAddRedirectBody struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func handleAddRedirect(c *gin.Context) {
	var body handleAddRedirectBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse body: %v", err))
		return
	}

	redirect := types.ProxyRedirect{
		Source: body.Source,
		Target: body.Target,
	}

	err = proxyService.AddRedirect(redirect)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func handleRemoveRedirect(c *gin.Context) {
	idString := c.Param("id")
	if idString == "" {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("failed to get redirection uuid"))
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = proxyService.RemoveRedirect(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
