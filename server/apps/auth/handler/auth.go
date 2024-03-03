package handler

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type authHandler struct {
	authService port.AuthService
}

func NewAuthHandler(authService port.AuthService) port.AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h authHandler) Login() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) (*types.Session, error) {
		login, pass, err := h.getUserPassFromHeader(ctx)
		if err != nil {
			return nil, err
		}
		token, err := h.authService.Login(login, pass)
		return &token, err
	})
}

func (h authHandler) Register() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) (*types.Session, error) {
		login, pass, err := h.getUserPassFromHeader(ctx)
		if err != nil {
			return nil, err
		}

		token, err := h.authService.Register(login, pass)
		return &token, err
	})
}

func (h authHandler) Verify() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) (*types.Session, error) {
		token := ctx.MustGet("token").(string)
		session, err := h.authService.Verify(token)
		return session, err
	})
}

func (h authHandler) Logout() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) error {
		token := ctx.MustGet("token").(string)
		return h.authService.Logout(token)
	})
}

func (h authHandler) getUserPassFromHeader(ctx *gin.Context) (string, string, error) {
	authorization := ctx.Request.Header.Get("Authorization")

	userpass := strings.TrimPrefix(authorization, "Basic ")
	userpassBytes, err := base64.StdEncoding.DecodeString(userpass)
	if err != nil {
		return "", "", err
	}
	userpass = string(userpassBytes)
	creds := strings.Split(userpass, ":")
	if len(creds) != 2 {
		return "", "", err
	}
	return creds[0], creds[1], nil
}
