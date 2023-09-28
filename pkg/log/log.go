package log

import (
	"os"
	"strings"

	"github.com/vertex-center/vlog"
)

var Default vlog.Logger
var DefaultKernel vlog.Logger

func init() {
	if strings.HasSuffix(os.Args[0], ".test") {
		Default = *vlog.New(
			vlog.WithOutputStd(),
		)
		DefaultKernel = *vlog.New(
			vlog.WithOutputStd(),
		)
		Default.Info("test logger initialized")
	} else {
		Default = *vlog.New(
			vlog.WithOutputStd(),
			vlog.WithOutputFile("live/logs", vlog.LogFormatText),
			vlog.WithOutputFile("live/logs", vlog.LogFormatJson),
		)
		DefaultKernel = *vlog.New(
			vlog.WithOutputStd(),
			vlog.WithOutputFile("live_kernel/logs", vlog.LogFormatText),
			vlog.WithOutputFile("live_kernel/logs", vlog.LogFormatJson),
		)
		Default.Info("full logger initialized")
	}
}

// Vertex

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

// Kernel

func DebugKernel(msg string, fields ...vlog.KeyValue) {
	Default.Debug(msg, fields...)
}

func InfoKernel(msg string, fields ...vlog.KeyValue) {
	Default.Info(msg, fields...)
}

func WarnKernel(msg string, fields ...vlog.KeyValue) {
	Default.Warn(msg, fields...)
}

func ErrorKernel(err error, fields ...vlog.KeyValue) {
	Default.Error(err, fields...)
}

func RequestKernel(msg string, fields ...vlog.KeyValue) {
	Default.Request(msg, fields...)
}
