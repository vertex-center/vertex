package api

import (
	"errors"
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
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
