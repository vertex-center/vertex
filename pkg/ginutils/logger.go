package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func Logger(router string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		if params.ErrorMessage != "" {
			log.Request("request",
				vlog.String("router", router),
				vlog.String("method", params.Method),
				vlog.Int("status", params.StatusCode),
				vlog.String("path", params.Path),
				vlog.String("latency", params.Latency.String()),
				vlog.String("ip", params.ClientIP),
				vlog.Int("size", params.BodySize),
				vlog.String("error", params.ErrorMessage),
			)
		} else {
			log.Request("request",
				vlog.String("router", router),
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
