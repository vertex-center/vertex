package server

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

func logger(u *url.URL) gin.HandlerFunc {
	urlString := u.String()
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		args := []vlog.KeyValue{
			vlog.String("request_id", params.Request.Header.Get("X-Request-ID")),
			vlog.String("correlation_id", params.Request.Header.Get("X-Correlation-ID")),
			vlog.String("url", urlString),
			vlog.String("method", params.Method),
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
	return func(c *gin.Context) {
		key := "X-Request-ID"
		id := uuid.New().String()

		c.Request.Header.Add(key, id)
		c.Header(key, id)
		c.Set("requestID", id)
		c.Next()
	}
}

// correlationID is a middleware that sets a unique ID
// for a request that is shared between services.
func correlationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "X-Correlation-ID"
		id := c.Request.Header.Get(key)
		if id != "" {
			// Check if id is valid, else ignore it.
			_, err := uuid.Parse(id)
			if err != nil {
				id = ""
			}
		}
		if id == "" {
			id = uuid.New().String()
			c.Request.Header.Add(key, id)
		}

		c.Header(key, id)
		c.Set("correlationID", id)
		c.Next()
	}
}
