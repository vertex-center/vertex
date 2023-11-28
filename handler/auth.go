package handler

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type AuthHandler struct {
	authService port.AuthService
}

func NewAuthHandler(authService port.AuthService) port.AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h AuthHandler) Login(c *router.Context) {
	c.OK()
}

func (h AuthHandler) Register(c *router.Context) {
	authorization := c.Request.Header.Get("Authorization")

	userpass := strings.TrimPrefix(authorization, "Basic ")
	userpassBytes, err := base64.StdEncoding.DecodeString(userpass)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           api.ErrInvalidCredentials,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: "Invalid credentials: expected base64 encoded login:password",
		})
		return
	}
	userpass = string(userpassBytes)
	creds := strings.Split(userpass, ":")
	if len(creds) != 2 {
		c.BadRequest(router.Error{
			Code:           api.ErrInvalidCredentials,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: "Invalid credentials: expected login:password",
		})
		return
	}

	err = h.authService.Register(creds[0], creds[1])
	if errors.Is(err, types.ErrLoginEmpty) {
		c.BadRequest(router.Error{
			Code:           api.ErrLoginEmpty,
			PublicMessage:  "Login must not be empty",
			PrivateMessage: "Login must not be empty",
		})
		return
	} else if errors.Is(err, types.ErrPasswordEmpty) {
		c.BadRequest(router.Error{
			Code:           api.ErrPasswordEmpty,
			PublicMessage:  "Password must not be empty",
			PrivateMessage: "Password must not be empty",
		})
		return
	} else if errors.Is(err, types.ErrPasswordLength) {
		c.BadRequest(router.Error{
			Code:           api.ErrPasswordLength,
			PublicMessage:  "Password must be at least 8 characters long",
			PrivateMessage: "Password must be at least 8 characters long",
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrInvalidCredentials,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: "Invalid credentials: " + err.Error(),
		})
		return
	}

	c.OK()
}
