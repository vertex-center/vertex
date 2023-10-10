package adapter

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
)

const bufferSize = 50

var (
	ErrLoggerNotFound = errors.New("instance logger not found")
)

type InstanceLogger struct {
	uuid uuid.UUID

	file        *os.File
	buffer      []instancestypes.LogLine
	currentLine int
	scheduler   *gocron.Scheduler

	dir string
}

type InstanceLogsFSAdapter struct {
	loggers      map[uuid.UUID]*InstanceLogger
	loggersMutex sync.RWMutex

	instancesPath string
}

type InstanceLogsFSAdapterParams struct {
	InstancesPath string
}

func NewInstanceLogsFSAdapter(params *InstanceLogsFSAdapterParams) instancestypes.InstanceLogsAdapterPort {
	if params == nil {
		params = &InstanceLogsFSAdapterParams{}
	}

	if params.InstancesPath == "" {
		params.InstancesPath = path.Join(storage.Path, "instances")
	}

	return &InstanceLogsFSAdapter{
		loggers:      map[uuid.UUID]*InstanceLogger{},
		loggersMutex: sync.RWMutex{},

		instancesPath: params.InstancesPath,
	}
}

func (a *InstanceLogsFSAdapter) Register(uuid uuid.UUID) error {
	dir := a.dir(uuid)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	l := InstanceLogger{
		uuid:   uuid,
		buffer: []instancestypes.LogLine{},
		dir:    dir,
	}

	a.loggersMutex.Lock()
	defer a.loggersMutex.Unlock()
	a.loggers[uuid] = &l

	err = l.Open()
	if err != nil {
		return err
	}

	return l.startCron()
}

func (a *InstanceLogsFSAdapter) Unregister(uuid uuid.UUID) error {
	l, err := a.getLogger(uuid)
	if err != nil {
		return err
	}

	err = l.stopCron()
	if err != nil {
		return err
	}

	err = l.Close()
	if err != nil {
		return err
	}

	a.loggersMutex.Lock()
	defer a.loggersMutex.Unlock()
	delete(a.loggers, uuid)
	return nil
}

func (a *InstanceLogsFSAdapter) Push(uuid uuid.UUID, line instancestypes.LogLine) {
	l, err := a.getLogger(uuid)
	if err != nil {
		log.Error(err)
		return
	}
	l.currentLine += 1
	l.buffer = append(l.buffer, line)
	if len(l.buffer) > bufferSize {
		l.buffer = l.buffer[1:]
	}

	_, err = fmt.Fprintf(l.file, "%s\n", line.Message.String())
	if err != nil {
		log.Error(err)
	}
}

func (a *InstanceLogsFSAdapter) Pop(uuid uuid.UUID) (instancestypes.LogLine, error) {
	l, err := a.getLogger(uuid)
	if err != nil {
		return instancestypes.LogLine{}, err
	}
	if len(l.buffer) == 0 {
		return instancestypes.LogLine{}, instancestypes.ErrBufferEmpty
	}
	line := l.buffer[len(l.buffer)-1]
	l.buffer = l.buffer[:len(l.buffer)-1]
	return line, nil
}

func (a *InstanceLogsFSAdapter) LoadBuffer(uuid uuid.UUID) ([]instancestypes.LogLine, error) {
	l, err := a.getLogger(uuid)
	if err != nil {
		return nil, err
	}
	return l.buffer, nil
}

func (a *InstanceLogsFSAdapter) UnregisterAll() error {
	for _, l := range a.loggers {
		err := a.Unregister(l.uuid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *InstanceLogsFSAdapter) getLogger(uuid uuid.UUID) (*InstanceLogger, error) {
	a.loggersMutex.RLock()
	defer a.loggersMutex.RUnlock()

	l, ok := a.loggers[uuid]
	if !ok {
		return nil, ErrLoggerNotFound
	}
	return l, nil
}

func (a *InstanceLogsFSAdapter) dir(uuid uuid.UUID) string {
	return path.Join(a.instancesPath, uuid.String(), ".vertex", "logs")
}

func (l *InstanceLogger) Open() error {
	filename := fmt.Sprintf("logs_%s.txt", time.Now().Format(time.DateOnly))
	filepath := path.Join(l.dir, filename)

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	l.file = file
	log.Info("opened instance logger", vlog.String("uuid", l.uuid.String()))
	return nil
}

func (l *InstanceLogger) Close() error {
	err := l.file.Close()
	if err != nil {
		return err
	}
	l.file = nil
	log.Info("closed instance logger", vlog.String("uuid", l.uuid.String()))
	return nil
}

func (l *InstanceLogger) startCron() error {
	l.scheduler = gocron.NewScheduler(time.Local)
	_, err := l.scheduler.Every(1).Day().At("00:00").Do(func() {
		err := l.Close()
		if err != nil {
			log.Error(err)
			return
		}
		err = l.Open()
		if err != nil {
			log.Error(err)
		}
	})
	if err != nil {
		return err
	}
	l.scheduler.StartAsync()
	return nil
}

func (l *InstanceLogger) stopCron() error {
	l.scheduler.Clear()
	l.scheduler.Stop()
	return nil
}
