package api

import (
	"errors"
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
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
