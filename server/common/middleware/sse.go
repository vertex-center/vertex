package middleware

import (
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

// SSE sets the headers for server-sent events.
func SSE(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", sse.ContentType)
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
