package handler

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
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

// docapi begin auth_login
// docapi method POST
// docapi summary Login
// docapi description Login with username and password
// docapi tags Authentication
// docapi response 200 {Token} The auth token
// docapi response 400
// docapi response 500
// docapi end

func (h AuthHandler) Login(c *router.Context) {
	login, pass, err := h.getUserPassFromHeader(c)
	if err != nil {
		return
	}

	token, err := h.authService.Login(login, pass)
	if errors.Is(err, types.ErrLoginFailed) {
		c.Abort(router.Error{
			Code:           types.ErrCodeInvalidCredentials,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(token)
}

// docapi begin auth_register
// docapi method POST
// docapi summary Register
// docapi description Register a new user with username and password
// docapi tags Authentication
// docapi response 200 {Token} The auth token
// docapi response 400
// docapi response 500

func (h AuthHandler) Register(c *router.Context) {
	login, pass, err := h.getUserPassFromHeader(c)
	if err != nil {
		return
	}

	token, err := h.authService.Register(login, pass)
	if errors.Is(err, types.ErrLoginEmpty) {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeLoginEmpty,
			PublicMessage:  "Login must not be empty",
			PrivateMessage: err.Error(),
		})
		return
	} else if errors.Is(err, types.ErrPasswordEmpty) {
		c.BadRequest(router.Error{
			Code:           types.ErrCodePasswordEmpty,
			PublicMessage:  "Password must not be empty",
			PrivateMessage: err.Error(),
		})
		return
	} else if errors.Is(err, types.ErrPasswordLength) {
		c.BadRequest(router.Error{
			Code:           types.ErrCodePasswordLength,
			PublicMessage:  "Password must be at least 8 characters long",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeInvalidCredentials,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(token)
}

// docapi begin auth_logout
// docapi method POST
// docapi summary Logout
// docapi tags Authentication
// docapi response 204
// docapi response 500
// docapi end

func (h AuthHandler) Logout(c *router.Context) {
	token := c.MustGet("token").(string)

	err := h.authService.Logout(token)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToLogout,
			PublicMessage:  "Failed to logout",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h AuthHandler) getUserPassFromHeader(c *router.Context) (string, string, error) {
	authorization := c.Request.Header.Get("Authorization")

	userpass := strings.TrimPrefix(authorization, "Basic ")
	userpassBytes, err := base64.StdEncoding.DecodeString(userpass)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeInvalidCredentials,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: "Expected base64 encoded login:password",
		})
		return "", "", err
	}
	userpass = string(userpassBytes)
	creds := strings.Split(userpass, ":")
	if len(creds) != 2 {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeInvalidCredentials,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: "Expected login:password",
		})
		return "", "", err
	}
	return creds[0], creds[1], nil
}
