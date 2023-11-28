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
	token := c.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	if AuthService == nil {
		log.Error(errors.New("auth_service is nil"))
		return
	}

	err := AuthService.Verify(token)
	if err != nil {
		c.Set("authenticated", false)
		c.Next()
		return
	}

	c.Set("authenticated", true)
	c.Set("token", token)
	c.Next()
}

func Authenticated(c *router.Context) {
	authenticated, exists := c.Get("authenticated")
	log.Debug("authenticated", vlog.Any("authenticated", authenticated))
	if !exists || !authenticated.(bool) {
		c.Abort(router.Error{
			Code:           types.ErrCodeInvalidToken,
			PublicMessage:  "Invalid token",
			PrivateMessage: "Invalid token",
		})
		return
	}
	c.Next()
}
