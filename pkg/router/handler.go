package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
)

func Handler(handler interface{}) gin.HandlerFunc {
	return HandleWithCode(handler, http.StatusOK)
}

func HandleWithCode(handler interface{}, code int) gin.HandlerFunc {
	return tonic.Handler(handler, code)
}
