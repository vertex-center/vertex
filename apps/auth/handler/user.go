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
// docapi tags Users
// docapi response 200 {User} The user
// docapi response 500
// docapi end

func (h *UserHandler) GetCurrentUser(c *router.Context) {
	userID := c.GetInt("user_id")

	user, err := h.service.GetUserByID(uint(userID))
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetUser,
			PublicMessage:  "Failed to retrieve the user profile",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(user)
}

// docapi begin auth_patch_current_user
// docapi method PATCH
// docapi summary Patch user
// docapi description Patch the logged-in user
// docapi tags Users
// docapi response 200 {User} The user
// docapi response 500
// docapi end

func (h *UserHandler) PatchCurrentUser(c *router.Context) {
	userID := c.GetInt("user_id")

	var user types.User
	err := c.ParseBody(&user)
	if err != nil {
		return
	}

	user.ID = uint(userID)
	user, err = h.service.PatchUser(user)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToPatchUser,
			PublicMessage:  "Failed to update the user profile",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(user)
}

// docapi begin auth_get_current_user_credentials
// docapi method GET
// docapi summary Get user credentials
// docapi description Retrieve the logged-in user credentials
// docapi tags Users
// docapi response 200 {UserCredentials} The user credentials
// docapi response 500
// docapi end

func (h *UserHandler) GetCurrentUserCredentials(c *router.Context) {
	userID := c.GetInt("user_id")

	credentials, err := h.service.GetUserCredentialsMethods(uint(userID))
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetUserCredentials,
			PublicMessage:  "Failed to retrieve the user credentials",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(credentials)
}
