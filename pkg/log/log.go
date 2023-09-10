package log

import "github.com/vertex-center/vlog"

var Default = *vlog.New(
	vlog.WithOutputStd(),
	vlog.WithOutputFile("live/logs", vlog.LogFormatText),
	vlog.WithOutputFile("live/logs", vlog.LogFormatJson),
)

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
