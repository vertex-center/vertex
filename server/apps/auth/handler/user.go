package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/apps/auth/core/types/session"
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
	return router.Handler(func(ctx *gin.Context) (*types.User, error) {
		s := session.Get(ctx)
		user, err := h.service.GetUserByID(s.UserID)
		return &user, err
	})
}

type PatchCurrentUserParams struct {
	types.User
}

func (h *userHandler) PatchCurrentUser() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *PatchCurrentUserParams) (*types.User, error) {
		s := session.Get(ctx)
		var err error
		params.ID = s.UserID
		params.User, err = h.service.PatchUser(params.User)
		if err != nil {
			return nil, err
		}
		return &params.User, nil
	})

}

func (h *userHandler) GetCurrentUserCredentials() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) ([]types.CredentialsMethods, error) {
		s := session.Get(ctx)
		return h.service.GetUserCredentialsMethods(s.UserID)
	})
}
