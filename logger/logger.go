package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	LogKindOut = "out"
	LogKindErr = "err"
)

var (
	tagInfo    = color.New(color.FgHiBlue).Sprintf("INF")
	tagWarn    = color.New(color.FgHiYellow).Sprintf("WRN")
	tagError   = color.New(color.FgHiRed).Sprintf("ERR")
	tagRequest = color.New(color.FgHiGreen).Sprintf("REQ")
)

type Line struct {
	tag     string
	kind    string
	color   color.Attribute
	message string
}

type S map[string]any

func formattedKeyValue(clr color.Attribute, key string, value any) string {
	message := color.New(clr).Sprintf("%s=", key)
	message += fmt.Sprintf("%v ", value)
	return message
}

func Log(message string) *Line {
	return &Line{
		tag:     tagInfo,
		kind:    LogKindOut,
		color:   color.FgBlue,
		message: formattedKeyValue(color.FgBlue, "msg", message),
	}
}

func Request() *Line {
	return &Line{
		tag:     tagRequest,
		kind:    LogKindOut,
		color:   color.FgGreen,
		message: "",
	}
}

func Warn(message string) *Line {
	return &Line{
		tag:     tagWarn,
		kind:    LogKindOut,
		color:   color.FgYellow,
		message: formattedKeyValue(color.FgYellow, "msg", message),
	}
}

func Error(err error) *Line {
	return &Line{
		tag:     tagError,
		kind:    LogKindErr,
		color:   color.FgRed,
		message: formattedKeyValue(color.FgRed, "msg", err.Error()),
	}
}

func (l *Line) AddKeyValue(key string, value any) *Line {
	l.message += formattedKeyValue(l.color, key, value)
	return l
}

func (l *Line) String() string {
	date := color.New(color.FgHiWhite).Sprintf(time.Now().Format(time.DateTime))
	return fmt.Sprintf("%s %s %s\n", date, l.tag, l.message)
}

func (l *Line) Print() {
	fmt.Print(l.String())
}
