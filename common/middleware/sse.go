package middleware

import (
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

// SSE sets the headers for server-sent events.
func SSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
