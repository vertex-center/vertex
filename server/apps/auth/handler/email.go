package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/apps/auth/core/types/session"
	"github.com/vertex-center/vertex/pkg/router"
)

type emailHandler struct {
	service port.EmailService
}

func NewEmailHandler(emailService port.EmailService) port.EmailHandler {
	return &emailHandler{
		service: emailService,
	}
}

func (h *emailHandler) GetCurrentUserEmails() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context) ([]types.Email, error) {
		s := session.Get(ctx)
		return h.service.GetEmails(s.UserID)
	})
}

type CreateCurrentUserEmailParams struct {
	Email string `json:"email"`
}

func (h *emailHandler) CreateCurrentUserEmail() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *CreateCurrentUserEmailParams) (*types.Email, error) {
		s := session.Get(ctx)
		email, err := h.service.CreateEmail(s.UserID, params.Email)
		return &email, err
	})
}

type DeleteCurrentUserEmailParams struct {
	Email string `json:"email"`
}

func (h *emailHandler) DeleteCurrentUserEmail() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *DeleteCurrentUserEmailParams) error {
		s := session.Get(ctx)
		return h.service.DeleteEmail(s.UserID, params.Email)
	})
}
