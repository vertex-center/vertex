package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type userHandler struct {
	service port.UserService
}

func NewUserHandler(userService port.UserService) port.UserHandler {
	return &userHandler{
		service: userService,
	}
}

func (h *userHandler) GetCurrentUser(c *router.Context) {
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

func (h *userHandler) GetCurrentUserInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get user"),
		oapi.Description("Retrieve the logged-in user"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.User{}),
		),
	}
}

func (h *userHandler) PatchCurrentUser(c *router.Context) {
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

func (h *userHandler) PatchCurrentUserInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Patch user"),
		oapi.Description("Update the logged-in user"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.User{}),
		),
	}
}

func (h *userHandler) GetCurrentUserCredentials(c *router.Context) {
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

func (h *userHandler) GetCurrentUserCredentialsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get user credentials"),
		oapi.Description("Retrieve the logged-in user credentials"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.CredentialsMethods{}),
		),
	}
}
