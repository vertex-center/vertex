package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/auth/api"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

func ReadAuth(ctx *gin.Context) {
	tokenStr := ctx.Request.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer")
	tokenStr = strings.TrimSpace(tokenStr)
	ctx.Set("token", tokenStr)
	log.Debug("reading auth", vlog.String("token", tokenStr))
	ctx.Next()
}

func Authenticated(ctx *gin.Context) {
	authClient := api.NewAuthClient(ctx)
	session, err := authClient.Verify(ctx)
	if err != nil {
		_ = ctx.AbortWithError(401, errors.NewUnauthorized(err, "invalid token"))
		return
	}

	ctx.Set("authenticated", true)
	ctx.Set("user_id", session.UserID)
	ctx.Next()
}
