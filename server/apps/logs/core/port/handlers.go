package port

import "github.com/gin-gonic/gin"

type LogsHandler interface {
	Push() gin.HandlerFunc
}
