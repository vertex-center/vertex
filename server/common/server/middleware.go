package server

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vlog"
)

const (
	HeaderXRequestID     = "X-Request-ID"
	HeaderXCorrelationID = "X-Correlation-ID"

	KeyRequestID     = "requestID"
	KeyCorrelationID = "correlationID"
)

func logger(u *url.URL, app string) gin.HandlerFunc {
	urlString := u.String()
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		args := []vlog.KeyValue{
			vlog.String("request_id", params.Request.Header.Get(HeaderXRequestID)),
			vlog.String("correlation_id", params.Request.Header.Get(HeaderXCorrelationID)),
			vlog.String("url", urlString),
			vlog.String("method", params.Method),
			vlog.String("app", app),
			vlog.Int("status", params.StatusCode),
			vlog.String("path", params.Path),
			vlog.String("latency", params.Latency.String()),
			vlog.String("ip", params.ClientIP),
			vlog.Int("size", params.BodySize),
		}
		if params.ErrorMessage != "" {
			errorMessage := strings.TrimSuffix(params.ErrorMessage, "\n")
			errorMessage = strings.ReplaceAll(errorMessage, "Error #01: ", "")
			args = append(args, vlog.String("error", errorMessage))
		}
		log.Request("request", args...)
		return ""
	})
}

// requestID is a middleware that sets a unique ID for each request.
func requestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := HeaderXRequestID
		id := uuid.New().String()

		ctx.Request.Header.Add(key, id)
		ctx.Header(key, id)
		ctx.Set(KeyRequestID, id)
		ctx.Next()
	}
}

// correlationID is a middleware that sets a unique ID
// for a request that is shared between services.
func correlationID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := HeaderXCorrelationID
		id := ctx.Request.Header.Get(key)
		if id != "" {
			// Check if id is valid, else ignore it.
			err := uuid.Validate(id)
			if err != nil {
				id = ""
			}
		}
		if id == "" {
			id = uuid.New().String()
			ctx.Request.Header.Add(key, id)
		}

		ctx.Header(key, id)
		ctx.Set(KeyCorrelationID, id)
		ctx.Next()
	}
}
