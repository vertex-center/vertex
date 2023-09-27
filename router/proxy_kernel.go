package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

func addProxyKernelRoutes(r *gin.RouterGroup) {
	r.POST("", handlePostProxyRedirects)
}

func handlePostProxyRedirects(c *gin.Context) {
	var redirects types.ProxyRedirects
	err := c.BindJSON(&redirects)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	proxyKernelService.SetRedirects(redirects)

	log.Info("proxy redirects updated",
		vlog.Int("count", len(redirects)),
	)

	c.Status(http.StatusNoContent)
}
