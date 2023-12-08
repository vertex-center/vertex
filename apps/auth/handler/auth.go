package handler

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
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

func (h authHandler) LoginInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("login"),
		fizz.Summary("Login"),
		fizz.Description("Login with username and password"),
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

func (h authHandler) RegisterInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("register"),
		fizz.Summary("Register"),
		fizz.Description("Register a new user with username and password"),
	}
}

func (h authHandler) Verify() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*types.Session, error) {
		token := c.MustGet("token").(string)
		session, err := h.authService.Verify(token)
		return session, err
	})
}

func (h authHandler) VerifyInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("verify"),
		fizz.Summary("Verify"),
		fizz.Description("Verify a token"),
	}
}

func (h authHandler) Logout() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) error {
		token := c.MustGet("token").(string)
		return h.authService.Logout(token)
	})
}

func (h authHandler) LogoutInfo() []fizz.OperationOption {
	return []fizz.OperationOption{
		fizz.ID("logout"),
		fizz.Summary("Logout"),
		fizz.Description("Logout a user"),
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
