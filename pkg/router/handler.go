package router

import "github.com/gin-gonic/gin"

type HandlerFunc func(*Context)

// wrapHandlers converts a slice of HandlerFunc to a slice of gin.HandlerFunc.
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	var ginHandlers []gin.HandlerFunc
	for _, handler := range handlers {
		// Capture the current value of handler. This is necessary because
		// otherwise the same handler will be used for each iteration.
		h := handler
		ginHandlers = append(ginHandlers, func(c *gin.Context) {
			h(&Context{c})
		})
	}
	return ginHandlers
}
