package handler

import (
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type emailHandler struct {
	service port.EmailService
}

func NewEmailHandler(emailService port.EmailService) port.EmailHandler {
	return &emailHandler{
		service: emailService,
	}
}

type GetCurrentUserEmailsParams struct {
	Email string `path:"email"`
}

func (h *emailHandler) GetCurrentUserEmails() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetCurrentUserEmailsParams) ([]types.Email, error) {
		return h.service.GetEmails(uint(c.GetInt("user_id")))
	})
}

func (h *emailHandler) GetCurrentUserEmailsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getCurrentUserEmails"),
		oapi.Summary("Get emails"),
		oapi.Description("Retrieve the emails of the logged-in user"),
	}
}

type CreateCurrentUserEmailParams struct {
	Email string `json:"email"`
}

func (h *emailHandler) CreateCurrentUserEmail() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *CreateCurrentUserEmailParams) (*types.Email, error) {
		userID := c.GetInt("user_id")

		addr, err := mail.ParseAddress(params.Email)
		if err != nil {
			return nil, err
		}

		email, err := h.service.CreateEmail(uint(userID), addr.Address)
		return &email, err
	})
}

func (h *emailHandler) CreateCurrentUserEmailInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("createCurrentUserEmail"),
		oapi.Summary("Create email"),
		oapi.Description("Create a new email for the logged-in user"),
	}
}

type DeleteCurrentUserEmailParams struct {
	Email string `path:"email"`
}

func (h *emailHandler) DeleteCurrentUserEmail() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *DeleteCurrentUserEmailParams) error {
		userID := c.GetInt("user_id")
		return h.service.DeleteEmail(uint(userID), params.Email)
	})
}

func (h *emailHandler) DeleteCurrentUserEmailInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("deleteCurrentUserEmail"),
		oapi.Summary("Delete email"),
		oapi.Description("Delete an email from the logged-in user"),
	}
}
