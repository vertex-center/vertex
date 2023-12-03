package handler

import (
	"errors"
	"net/mail"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type EmailHandler struct {
	service port.EmailService
}

func NewEmailHandler(emailService port.EmailService) port.EmailHandler {
	return &EmailHandler{
		service: emailService,
	}
}

// docapi begin auth_current_user_get_emails
// docapi method GET
// docapi summary Get emails
// docapi description Retrieve the emails of the logged-in user
// docapi tags Authentication
// docapi response 200 {[]Email} The emails
// docapi response 500
// docapi end

func (h *EmailHandler) GetCurrentUserEmails(c *router.Context) {
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

// docapi begin auth_current_user_create_email
// docapi method POST
// docapi summary Create email
// docapi description Create a new email for the logged-in user
// docapi tags Authentication
// docapi response 200 {Email} The email
// docapi response 400
// docapi response 409 {Error} Email already exists
// docapi response 500
// docapi end

type CreateCurrentUserEmailBody struct {
	Email string `json:"email"`
}

func (h *EmailHandler) CreateCurrentUserEmail(c *router.Context) {
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

// docapi begin auth_current_user_delete_email
// docapi method DELETE
// docapi summary Delete email
// docapi description Delete an email from the logged-in user
// docapi tags Authentication
// docapi response 204
// docapi response 500
// docapi end

func (h *EmailHandler) DeleteCurrentUserEmail(c *router.Context) {
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
