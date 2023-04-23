package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

var (
	DefaultLogger Logger

	LogKindOut = "out"
	LogKindErr = "err"

	tagInfo    = "INF"
	tagWarn    = "WRN"
	tagError   = "ERR"
	tagRequest = "REQ"
)

type Logger struct {
	out *os.File
	err *os.File

	text *os.File
	json *os.File
}

type Line struct {
	logger         *Logger
	tag            string
	kind           string
	color          color.Attribute
	date           string
	messageColored string
	messagePlain   string
	json           map[string]any
}

func NewDefaultLogger() Logger {
	t, err := os.OpenFile("vertex_logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open vertex.logs: %v", err)
		t = nil
	}

	// jsonl stands for the json lines format. https://jsonlines.org/
	j, err := os.OpenFile("vertex_logs.jsonl", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open vertex.logs: %v", err)
		j = nil
	}

	return Logger{
		out:  os.Stdout,
		err:  os.Stderr,
		text: t,
		json: j,
	}
}

func (l *Logger) Close() {
	l.json.Close()
}

func (l *Logger) Date() string {
	return time.Now().Format(time.DateTime)
}

func Log(message string) *Line {
	return DefaultLogger.Log(message)
}

func Request() *Line {
	return DefaultLogger.Request()
}

func Warn(message string) *Line {
	return DefaultLogger.Warn(message)
}

func Error(err error) *Line {
	return DefaultLogger.Error(err)
}

func (l *Logger) Log(message string) *Line {
	return &Line{
		logger:         l,
		tag:            tagInfo,
		kind:           LogKindOut,
		color:          color.FgBlue,
		date:           l.Date(),
		messageColored: formattedKeyValue(color.FgBlue, "msg", message),
		messagePlain:   "msg=" + message,
		json: map[string]any{
			"seconds":     time.Now().Unix(),
			"nanoseconds": time.Now().UnixNano(),
			"kind":        "info",
			"msg":         message,
		},
	}
}

func (l *Logger) Request() *Line {
	return &Line{
		logger:         l,
		tag:            tagRequest,
		kind:           LogKindOut,
		color:          color.FgGreen,
		date:           l.Date(),
		messageColored: "",
		messagePlain:   "",
		json: map[string]any{
			"seconds":     time.Now().Unix(),
			"nanoseconds": time.Now().UnixNano(),
			"kind":        "request",
		},
	}
}

func (l *Logger) Warn(message string) *Line {
	return &Line{
		logger:         l,
		tag:            tagWarn,
		kind:           LogKindOut,
		color:          color.FgYellow,
		date:           l.Date(),
		messageColored: formattedKeyValue(color.FgYellow, "msg", message),
		messagePlain:   "msg=" + message,
		json: map[string]any{
			"seconds":     time.Now().Unix(),
			"nanoseconds": time.Now().UnixNano(),
			"kind":        "warn",
			"msg":         message,
		},
	}
}

func (l *Logger) Error(err error) *Line {
	return &Line{
		logger:         l,
		tag:            tagError,
		kind:           LogKindErr,
		color:          color.FgRed,
		date:           l.Date(),
		messageColored: formattedKeyValue(color.FgRed, "msg", err.Error()),
		messagePlain:   "msg=" + err.Error(),
		json: map[string]any{
			"seconds":     time.Now().Unix(),
			"nanoseconds": time.Now().UnixNano(),
			"kind":        "error",
			"msg":         err.Error(),
		},
	}
}

func (l *Line) AddKeyValue(key string, value any) *Line {
	l.messageColored += formattedKeyValue(l.color, key, value)
	l.messagePlain += fmt.Sprintf("%s=%v", key, value)
	l.json[key] = value
	return l
}

func (l *Line) String() string {
	return fmt.Sprintf("%s %s %s\n",
		color.New(color.FgHiWhite).Sprintf(l.date),
		color.New(l.color).Sprintf(l.tag),
		l.messageColored,
	)
}

func (l *Line) StringPlain() string {
	return fmt.Sprintf("%s %s %s\n", l.date, l.tag, l.messagePlain)
}

func (l *Line) Json() string {
	j, err := json.Marshal(l.json)
	if err != nil {
		return ""
	}
	return string(j) + "\n"
}

func (l *Line) Print() {
	if l.kind == LogKindErr {
		_, _ = fmt.Fprint(l.logger.err, l.String())
	} else {
		_, _ = fmt.Fprint(l.logger.out, l.String())
	}
	if l.logger.text != nil {
		_, _ = fmt.Fprint(l.logger.text, l.StringPlain())
	}
	if l.logger.json != nil {
		_, _ = fmt.Fprint(l.logger.json, l.Json())
	}
}

func formattedKeyValue(clr color.Attribute, key string, value any) string {
	message := color.New(clr).Sprintf("%s=", key)
	message += fmt.Sprintf("%v ", value)

	return message
}
