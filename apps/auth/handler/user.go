package handler

import (
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
	userID := c.GetInt("user_id")

	user, err := h.service.GetUserByID(uint(userID))
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetUser,
			PublicMessage:  "Failed to retrieve the user",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(user)
}
