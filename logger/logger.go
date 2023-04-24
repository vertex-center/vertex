package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/fatih/color"
	"github.com/go-co-op/gocron"
)

var (
	DefaultLogger *Logger

	LogKindOut = "out"
	LogKindErr = "err"

	tagInfo    = "INF"
	tagWarn    = "WRN"
	tagError   = "ERR"
	tagRequest = "REQ"

	logsPath = path.Join("live", "logs")
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

func NewDefaultLogger() *Logger {
	_ = os.MkdirAll(logsPath, os.ModePerm)

	l := &Logger{
		out: os.Stdout,
		err: os.Stderr,
	}

	l.OpenLogFiles()
	l.StartCron()

	return l
}

func (l *Logger) StartCron() {
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(1).Day().At("00:00").Do(func() {
		l.Close()
		l.OpenLogFiles()
	})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
	}
	s.StartAsync()
}

func (l *Logger) Close() {
	l.text.Close()
	l.json.Close()
}

func (l *Logger) OpenLogFiles() {
	filename := fmt.Sprintf("vertex_logs_%s.txt", time.Now().Format(time.DateOnly))
	t, err := os.OpenFile(path.Join(logsPath, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", filename, err)
		t = nil
	}
	l.text = t

	// jsonl stands for the json lines format. https://jsonlines.org/
	filename = fmt.Sprintf("vertex_logs_%s.jsonl", time.Now().Format(time.DateOnly))
	j, err := os.OpenFile(path.Join(logsPath, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", filename, err)
		j = nil
	}
	l.json = j
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

func (l *Logger) NewLine(tag string, kind string, color color.Attribute, message string) *Line {
	messageColored := ""
	messagePlain := ""

	if message != "" {
		messageColored = formattedKeyValue(color, "msg", message)
		messagePlain = "msg=" + message + " "
	}

	return &Line{
		logger:         l,
		tag:            tag,
		kind:           kind,
		color:          color,
		date:           l.Date(),
		messageColored: messageColored,
		messagePlain:   messagePlain,
		json: map[string]any{
			"seconds":     time.Now().Unix(),
			"nanoseconds": time.Now().UnixNano(),
			"kind":        tag,
			"msg":         message,
		},
	}
}

func (l *Logger) Log(message string) *Line {
	return l.NewLine(tagInfo, LogKindOut, color.FgBlue, message)
}

func (l *Logger) Request() *Line {
	return l.NewLine(tagRequest, LogKindOut, color.FgGreen, "")
}

func (l *Logger) Warn(message string) *Line {
	return l.NewLine(tagWarn, LogKindOut, color.FgYellow, message)
}

func (l *Logger) Error(err error) *Line {
	return l.NewLine(tagError, LogKindErr, color.FgRed, err.Error())
}

func (l *Line) AddKeyValue(key string, value any) *Line {
	l.messageColored += formattedKeyValue(l.color, key, value)
	l.messagePlain += fmt.Sprintf("%s=%v ", key, value)
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
	l.PrintInExternalFiles()
}

func (l *Line) PrintInExternalFiles() {
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
