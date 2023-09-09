package log

import "github.com/vertex-center/vlog"

var Default = *vlog.New(
	vlog.WithOutputStd(),
	vlog.WithOutputFile("live/logs", vlog.LogFormatText),
	vlog.WithOutputFile("live/logs", vlog.LogFormatJson),
)
