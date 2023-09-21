package ginutils

import "github.com/gin-gonic/gin"

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(-1, c.Errors[0].Err)
		}
	}
}
