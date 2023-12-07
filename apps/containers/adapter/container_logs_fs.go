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
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

const bufferSize = 50

var (
	ErrLoggerNotFound = errors.New("container logger not found")
)

type ContainerLogger struct {
	uuid uuid.UUID

	file        *os.File
	buffer      []types.LogLine
	currentLine int
	scheduler   *gocron.Scheduler

	dir string
}

type containerLogsFSAdapter struct {
	loggers      map[uuid.UUID]*ContainerLogger
	loggersMutex sync.RWMutex

	containersPath string
}

type ContainerLogsFSAdapterParams struct {
	ContainersPath string
}

func NewContainerLogsFSAdapter(params *ContainerLogsFSAdapterParams) port.ContainerLogsAdapter {
	if params == nil {
		params = &ContainerLogsFSAdapterParams{}
	}

	if params.ContainersPath == "" {
		params.ContainersPath = path.Join(storage.FSPath, "apps", "containers", "containers")
	}

	return &containerLogsFSAdapter{
		loggers:      map[uuid.UUID]*ContainerLogger{},
		loggersMutex: sync.RWMutex{},

		containersPath: params.ContainersPath,
	}
}

func (a *containerLogsFSAdapter) Register(uuid uuid.UUID) error {
	dir := a.dir(uuid)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	l := ContainerLogger{
		uuid:   uuid,
		buffer: []types.LogLine{},
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

func (a *containerLogsFSAdapter) Unregister(uuid uuid.UUID) error {
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

func (a *containerLogsFSAdapter) Push(uuid uuid.UUID, line types.LogLine) {
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

func (a *containerLogsFSAdapter) Pop(uuid uuid.UUID) (types.LogLine, error) {
	l, err := a.getLogger(uuid)
	if err != nil {
		return types.LogLine{}, err
	}
	if len(l.buffer) == 0 {
		return types.LogLine{}, types.ErrBufferEmpty
	}
	line := l.buffer[len(l.buffer)-1]
	l.buffer = l.buffer[:len(l.buffer)-1]
	return line, nil
}

func (a *containerLogsFSAdapter) LoadBuffer(uuid uuid.UUID) ([]types.LogLine, error) {
	l, err := a.getLogger(uuid)
	if err != nil {
		return nil, err
	}
	return l.buffer, nil
}

func (a *containerLogsFSAdapter) UnregisterAll() error {
	var ids []uuid.UUID

	a.loggersMutex.RLock()
	for id := range a.loggers {
		ids = append(ids, id)
	}
	a.loggersMutex.RUnlock()

	for _, id := range ids {
		log.Info("unregistering container logger", vlog.String("uuid", id.String()))
		err := a.Unregister(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *containerLogsFSAdapter) getLogger(uuid uuid.UUID) (*ContainerLogger, error) {
	a.loggersMutex.RLock()
	defer a.loggersMutex.RUnlock()

	l, ok := a.loggers[uuid]
	if !ok {
		return nil, ErrLoggerNotFound
	}
	return l, nil
}

func (a *containerLogsFSAdapter) dir(uuid uuid.UUID) string {
	return path.Join(a.containersPath, uuid.String(), ".vertex", "logs")
}

func (l *ContainerLogger) Open() error {
	filename := fmt.Sprintf("logs_%s.txt", time.Now().Format(time.DateOnly))
	filepath := path.Join(l.dir, filename)

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	l.file = file
	log.Info("opened container logger", vlog.String("uuid", l.uuid.String()))
	return nil
}

func (l *ContainerLogger) Close() error {
	err := l.file.Close()
	if err != nil {
		return err
	}
	l.file = nil
	log.Info("closed container logger", vlog.String("uuid", l.uuid.String()))
	return nil
}

func (l *ContainerLogger) startCron() error {
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

func (l *ContainerLogger) stopCron() error {
	l.scheduler.Clear()
	l.scheduler.Stop()
	return nil
}
