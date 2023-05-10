package types

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/vertex-center/vertex/pkg/logger"
)

type InstanceLogger struct {
	file *os.File

	currentLine int
	logsDir     string
}

func NewInstanceLogger(instancePath string) *InstanceLogger {
	dir := path.Join(instancePath, ".vertex", "logs")

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.Error(err).Print()
	}

	l := &InstanceLogger{
		logsDir: dir,
	}
	l.OpenLogFile()
	l.StartCron()
	return l
}

func (l *InstanceLogger) OpenLogFile() {
	filename := fmt.Sprintf("logs_%s.txt", time.Now().Format(time.DateOnly))
	filepath := path.Join(l.logsDir, filename)

	var err error
	l.file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		logger.Error(err).Print()
	}
}

func (l *InstanceLogger) CloseLogFile() {
	l.file.Close()
}

func (l *InstanceLogger) StartCron() {
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(1).Day().At("00:00").Do(func() {
		l.CloseLogFile()
		l.OpenLogFile()
	})
	if err != nil {
		logger.Error(err).Print()
		return
	}
	s.StartAsync()
}

func (l *InstanceLogger) Write(line *LogLine) {
	l.currentLine += 1
	line.Id = l.currentLine

	_, err := fmt.Fprintf(l.file, "%s\n", line.Message)
	if err != nil {
		logger.Error(err).Print()
	}
}
