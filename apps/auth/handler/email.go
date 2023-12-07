package handler

import (
	"errors"
	"net/http"
	"net/mail"

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

func (h *emailHandler) GetCurrentUserEmails(c *router.Context) {
	userID := c.GetInt("user_id")

	emails, err := h.service.GetEmails(uint(userID))
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToGetUserEmails,
			PublicMessage:  "Failed to retrieve your email addresses",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(emails)
}

func (h *emailHandler) GetCurrentUserEmailsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get emails"),
		oapi.Description("Retrieve the emails of the logged-in user"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel([]types.Email{}),
		),
	}
}

type CreateCurrentUserEmailBody struct {
	Email string `json:"email"`
}

func (h *emailHandler) CreateCurrentUserEmail(c *router.Context) {
	userID := c.GetInt("user_id")

	var body CreateCurrentUserEmailBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	addr, err := mail.ParseAddress(body.Email)
	if err != nil {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeInvalidEmail,
			PublicMessage:  "This email address is not a valid email address",
			PrivateMessage: err.Error(),
		})
		return
	}

	email, err := h.service.CreateEmail(uint(userID), addr.Address)
	if errors.Is(err, types.ErrEmailAlreadyExists) {
		c.Conflict(router.Error{
			Code:           types.ErrCodeEmailAlreadyExists,
			PublicMessage:  "This email address is already registered on your account",
			PrivateMessage: err.Error(),
		})
		return
	} else if errors.Is(err, types.ErrEmailEmpty) {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeEmailEmpty,
			PublicMessage:  "Email address must not be empty",
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToCreateEmail,
			PublicMessage:  "Failed to add this email address",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(email)
}

func (h *emailHandler) CreateCurrentUserEmailInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Create email"),
		oapi.Description("Create a new email for the logged-in user"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Email{}),
		),
	}
}

func (h *emailHandler) DeleteCurrentUserEmail(c *router.Context) {
	userID := c.GetInt("user_id")

	var body CreateCurrentUserEmailBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	err = h.service.DeleteEmail(uint(userID), body.Email)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToDeleteEmail,
			PublicMessage:  "Failed to delete this email address",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *emailHandler) DeleteCurrentUserEmailInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Delete email"),
		oapi.Description("Delete an email from the logged-in user"),
		oapi.Response(http.StatusNoContent),
	}
}
