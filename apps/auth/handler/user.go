package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type userHandler struct {
	service port.UserService
}

func NewUserHandler(userService port.UserService) port.UserHandler {
	return &userHandler{
		service: userService,
	}
}

func (h *userHandler) GetCurrentUser() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*types.User, error) {
		userID := c.GetInt("user_id")
		user, err := h.service.GetUserByID(uint(userID))
		return &user, err
	})
}

type PatchCurrentUserParams struct {
	types.User
}

func (h *userHandler) PatchCurrentUser() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *PatchCurrentUserParams) (*types.User, error) {
		userID := c.GetInt("user_id")
		var err error
		params.ID = uint(userID)
		params.User, err = h.service.PatchUser(params.User)
		if err != nil {
			return nil, err
		}
		return &params.User, nil
	})

}

func (h *userHandler) GetCurrentUserCredentials() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.CredentialsMethods, error) {
		userID := c.GetInt("user_id")
		return h.service.GetUserCredentialsMethods(uint(userID))
	})
}
