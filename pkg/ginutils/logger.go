package ginutils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func Logger(router, port string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		if params.ErrorMessage != "" {
			errorMessage := strings.TrimSuffix(params.ErrorMessage, "\n")
			errorMessage = strings.ReplaceAll(errorMessage, "Error #01: ", "")

			log.Request("request",
				vlog.String("router", router),
				vlog.String("port", port),
				vlog.String("method", params.Method),
				vlog.Int("status", params.StatusCode),
				vlog.String("path", params.Path),
				vlog.String("latency", params.Latency.String()),
				vlog.String("ip", params.ClientIP),
				vlog.Int("size", params.BodySize),
				vlog.String("error", errorMessage),
			)
		} else {
			log.Request("request",
				vlog.String("router", router),
				vlog.String("port", port),
				vlog.String("method", params.Method),
				vlog.Int("status", params.StatusCode),
				vlog.String("path", params.Path),
				vlog.String("latency", params.Latency.String()),
				vlog.String("ip", params.ClientIP),
				vlog.Int("size", params.BodySize),
			)
		}

		return ""
	})
}
