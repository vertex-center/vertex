package log

import (
	"fmt"
	"os"
	"strings"

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
				err := defaultAgent.Send(line)
				if err != nil {
					_, _ = fmt.Fprint(os.Stderr, err.Error())
				}
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
