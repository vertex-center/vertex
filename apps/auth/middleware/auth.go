package middleware

import (
	"errors"
	"strings"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

var AuthService port.AuthService

func ReadAuth(c *router.Context) {
	tokenStr := c.Request.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer")
	tokenStr = strings.TrimSpace(tokenStr)

	if AuthService == nil {
		log.Error(errors.New("auth_service is nil"))
		c.Set("authenticated", false)
		c.Next()
		return
	}

	log.Debug("reading auth", vlog.String("token", tokenStr))

	session, err := AuthService.Verify(tokenStr)
	if err != nil {
		c.Set("authenticated", false)
		c.Next()
		return
	}

	c.Set("authenticated", true)
	c.Set("token", tokenStr)
	c.Set("user_id", int(session.UserID))
	c.Next()
}

func Authenticated(c *router.Context) {
	authenticated, exists := c.Get("authenticated")
	log.Debug("authenticated", vlog.Any("authenticated", authenticated))
	if !exists || !authenticated.(bool) {
		c.Unauthorized(router.Error{
			Code:           types.ErrCodeInvalidToken,
			PublicMessage:  "Invalid token",
			PrivateMessage: "Invalid token",
		})
		return
	}
	c.Next()
}
