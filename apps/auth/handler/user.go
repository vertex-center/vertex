package handler

import (
	"fmt"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(userService port.UserService) port.UserHandler {
	return &UserHandler{
		service: userService,
	}
}

// docapi begin auth_get_current_user
// docapi method GET
// docapi summary Get user
// docapi description Retrieve the logged-in user
// docapi tags Authentication/Users
// docapi response 200 {User} The user
// docapi response 500
// docapi end

func (h *UserHandler) GetCurrentUser(c *router.Context) {
	username := c.GetString("username")

	user, err := h.service.GetUser(username)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetUser,
			PublicMessage:  fmt.Sprintf("Failed to retrieve the user '%s'", username),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(user)
}
