package log

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/vertex-center/vertex/apps/logs/api"
	"github.com/vertex-center/vlog"
)

var Default vlog.Logger

func init() {
	if strings.HasSuffix(os.Args[0], ".test") {
		Default = *vlog.New(
			vlog.WithOutputStd(),
		)
		Default.Info("test logger initialized")
	} else {
		Default = *vlog.New(
			vlog.WithOutputStd(),
			vlog.WithOutputFunc(vlog.LogFormatJson, func(line string) {
				go func() {
					logsClient := api.NewLogsClient()
					err := logsClient.PushLogs(context.Background(), line)
					if err != nil {
						fmt.Println("failed to push logs:", err)
					}
				}()
			}),
		)
		Default.Info("full logger initialized")
	}
}

func Debug(msg string, fields ...vlog.KeyValue) {
	Default.Debug(msg, fields...)
}

func Info(msg string, fields ...vlog.KeyValue) {
	Default.Info(msg, fields...)
}

func Warn(msg string, fields ...vlog.KeyValue) {
	Default.Warn(msg, fields...)
}

func Error(err error, fields ...vlog.KeyValue) {
	Default.Error(err, fields...)
}

func Request(msg string, fields ...vlog.KeyValue) {
	Default.Request(msg, fields...)
}
