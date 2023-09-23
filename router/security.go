package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addSecurityRoutes(r *gin.RouterGroup) {
	r.GET("/ssh", handleGetSSHKey)
}

// handleGetSSHKey handles the retrieval of the SSH key.
func handleGetSSHKey(c *gin.Context) {
	keys, err := sshService.GetAll()
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}

	c.JSON(http.StatusOK, keys)
}
