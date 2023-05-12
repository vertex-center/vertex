package repository

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
)

type InstanceLogger struct {
	file *os.File

	currentLine int
}

type InstanceLogsFSRepository struct {
	loggers map[uuid.UUID]*InstanceLogger
}

func NewInstanceLogsFSRepository() InstanceLogsFSRepository {
	r := InstanceLogsFSRepository{
		loggers: map[uuid.UUID]*InstanceLogger{},
	}
	r.startCron()
	return r
}

func (r *InstanceLogsFSRepository) Open(uuid uuid.UUID) error {
	dir := path.Join(storage.PathInstances, uuid.String(), ".vertex", "logs")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("logs_%s.txt", time.Now().Format(time.DateOnly))
	filepath := path.Join(storage.PathInstances, uuid.String(), ".vertex", "logs", filename)

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	l := InstanceLogger{}
	l.file = file

	r.loggers[uuid] = &l
	return nil
}

func (r *InstanceLogsFSRepository) Close(uuid uuid.UUID) error {
	l := r.loggers[uuid]
	return l.Close()
}

func (r *InstanceLogsFSRepository) Push(uuid uuid.UUID, line types.LogLine) {
	l := r.loggers[uuid]
	l.currentLine += 1

	_, err := fmt.Fprintf(l.file, "%s\n", line.Message)
	if err != nil {
		logger.Error(err).Print()
	}
}

func (r *InstanceLogsFSRepository) CloseAll() error {
	var errs []error
	for id := range r.loggers {
		err := r.Close(id)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (l *InstanceLogger) Close() error {
	return l.file.Close()
}

func (r *InstanceLogsFSRepository) startCron() {
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(1).Day().At("00:00").Do(func() {
		for id := range r.loggers {
			err := r.Close(id)
			if err != nil {
				logger.Error(err).Print()
				continue
			}
			err = r.Open(id)
			if err != nil {
				logger.Error(err).Print()
			}
		}
	})
	if err != nil {
		logger.Error(err).Print()
		return
	}
	s.StartAsync()
}
