package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

type Error struct {
	HttpCode int
	Code     router.ErrCode `json:"code"`
	Message  string         `json:"message"`
}

func (e *Error) RouterError() router.Error {
	return router.Error{
		Code:          e.Code,
		PublicMessage: e.Message,
	}
}

func HandleError(requestError error, apiError Error) *Error {
	if errors.Is(requestError, requests.ErrValidator) {
		if requests.HasStatusErr(requestError, http.StatusNotFound) {
			apiError.HttpCode = http.StatusNotFound
			apiError.Message = "Resource not found."
		} else if requests.HasStatusErr(requestError, http.StatusInternalServerError) {
			apiError.HttpCode = http.StatusInternalServerError
		}
		return &apiError
	} else if requestError != nil {
		log.Error(requestError)
		return &Error{
			HttpCode: http.StatusInternalServerError,
			Code:     ErrInternalError,
			Message:  "Internal error.",
		}
	}
	return nil
}

func AuthMiddleware(c *router.Context, service port.AuthService) {
	token := c.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		c.BadRequest(router.Error{
			Code:           ErrAuthorizationEmpty,
			PublicMessage:  "Authorization header is empty",
			PrivateMessage: "Authorization header is empty",
		})
		return
	}

	err := service.Verify(token)
	if err != nil {
		c.Abort(router.Error{
			Code:           ErrInvalidToken,
			PublicMessage:  "Invalid credentials",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.Set("token", token)
	c.Next()
}

func AppRequest(AppRoute string) *requests.Builder {
	return requests.New(func(rb *requests.Builder) {
		path := "/api/app" + AppRoute + "/"
		rb.BaseURL(config.Current.VertexURL()).Path(path)
		rb.AddValidator(func(response *http.Response) error {
			log.Request("request from app", vlog.String("path", response.Request.URL.Path))
			return nil
		})
	})
}
