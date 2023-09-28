package log

import (
	"os"
	"strings"

	"github.com/vertex-center/vlog"
)

var Default vlog.Logger

func init() {
	var p string
	if strings.Contains(os.Args[0], "kernel") || strings.Contains(os.Args[0], "Kernel") {
		p = "live_kernel/logs"
	} else {
		p = "live/logs"
	}

	if strings.HasSuffix(os.Args[0], ".test") {
		Default = *vlog.New(
			vlog.WithOutputStd(),
		)
		Default.Info("test logger initialized")
	} else {
		Default = *vlog.New(
			vlog.WithOutputStd(),
			vlog.WithOutputFile(p, vlog.LogFormatText),
			vlog.WithOutputFile(p, vlog.LogFormatJson),
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
