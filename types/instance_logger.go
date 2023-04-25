package types

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/vertex-center/vertex/logger"
)

type InstanceLogger struct {
	file *os.File

	currentLine int
}

func NewInstanceLogger(instancePath string) (*InstanceLogger, error) {
	dir := path.Join(instancePath, ".vertex", "logs")

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("logs_%s.txt", time.Now().Format(time.DateOnly))
	filepath := path.Join(dir, filename)

	l := &InstanceLogger{}

	l.file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return l, err
}

func (l *InstanceLogger) Close() {
	l.file.Close()
}

func (l *InstanceLogger) Write(line *LogLine) {
	l.currentLine += 1
	line.Id = l.currentLine

	_, err := fmt.Fprintf(l.file, "%s\n", line.Message)
	if err != nil {
		logger.Error(err).Print()
	}
}
