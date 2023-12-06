package middleware

import (
	"strings"

	"github.com/vertex-center/vertex/apps/auth/api"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

func ReadAuth(c *router.Context) {
	tokenStr := c.Request.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer")
	tokenStr = strings.TrimSpace(tokenStr)
	c.Set("token", tokenStr)
	log.Debug("reading auth", vlog.String("token", tokenStr))
}

func Authenticated(c *router.Context) {
	tokenStr := c.GetString("token")

	authClient := api.NewAuthClient(tokenStr)
	session, err := authClient.Verify(c)
	if err != nil {
		c.Unauthorized(router.Error{
			Code:           types.ErrCodeInvalidToken,
			PublicMessage:  "Invalid token",
			PrivateMessage: "Invalid token",
		})
		return
	}

	c.Set("authenticated", true)
	c.Set("user_id", int(session.UserID))
	c.Next()
}
