package ginutils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/logger"
)

func Logger(router string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		l := logger.Request().
			AddKeyValue("router", router).
			AddKeyValue("method", params.Method).
			AddKeyValue("status", params.StatusCode).
			AddKeyValue("path", params.Path).
			AddKeyValue("latency", params.Latency).
			AddKeyValue("ip", params.ClientIP).
			AddKeyValue("size", params.BodySize)

		if params.ErrorMessage != "" {
			err, _ := strings.CutSuffix(params.ErrorMessage, "\n")
			l.AddKeyValue("error", err)
		}

		l.PrintInExternalFiles()

		return l.String()
	})
}
