package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/auth/api"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

func ReadAuth(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer")
	tokenStr = strings.TrimSpace(tokenStr)
	c.Set("token", tokenStr)
	log.Debug("reading auth", vlog.String("token", tokenStr))
}

func Authenticated(c *gin.Context) {
	tokenStr := c.GetString("token")

	authClient := api.NewAuthClient(tokenStr)
	session, err := authClient.Verify(c)
	if err != nil {
		_ = c.AbortWithError(401, errors.NewUnauthorized(err, "invalid token"))
		return
	}

	c.Set("authenticated", true)
	c.Set("user_id", int(session.UserID))
	c.Next()
}
