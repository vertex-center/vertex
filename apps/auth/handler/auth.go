package handler

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
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
	return router.Handler(func(c *gin.Context) (*types.Session, error) {
		login, pass, err := h.getUserPassFromHeader(c)
		if err != nil {
			return nil, err
		}
		token, err := h.authService.Login(login, pass)
		return &token, err
	})
}

func (h authHandler) LoginInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("login"),
		oapi.Summary("Login"),
		oapi.Description("Login with username and password"),
	}
}

func (h authHandler) Register() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*types.Session, error) {
		login, pass, err := h.getUserPassFromHeader(c)
		if err != nil {
			return nil, err
		}

		token, err := h.authService.Register(login, pass)
		return &token, err
	})
}

func (h authHandler) RegisterInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("register"),
		oapi.Summary("Register"),
		oapi.Description("Register a new user with username and password"),
	}
}

func (h authHandler) Verify() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*types.Session, error) {
		token := c.MustGet("token").(string)
		session, err := h.authService.Verify(token)
		return session, err
	})
}

func (h authHandler) VerifyInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("verify"),
		oapi.Summary("Verify"),
		oapi.Description("Verify a token"),
	}
}

func (h authHandler) Logout() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) error {
		token := c.MustGet("token").(string)
		return h.authService.Logout(token)
	})
}

func (h authHandler) LogoutInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("logout"),
		oapi.Summary("Logout"),
		oapi.Description("Logout a user"),
	}
}

func (h authHandler) getUserPassFromHeader(c *gin.Context) (string, string, error) {
	authorization := c.Request.Header.Get("Authorization")

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
