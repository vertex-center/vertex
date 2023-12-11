package adapter

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vlog"
)

const bufferSize = 50

var (
	ErrLoggerNotFound = errors.NotFoundf("logger")
)

type ContainerLogger struct {
	uuid        types.ContainerID
	file        *os.File
	buffer      []types.LogLine
	currentLine int
	scheduler   *gocron.Scheduler
	dir         string
}

type logsFSAdapter struct {
	loggers   map[types.ContainerID]*ContainerLogger
	loggersMu sync.RWMutex

	containersPath string
}

type LogsFSAdapterParams struct {
	ContainersPath string
}

func NewLogsFSAdapter(params *LogsFSAdapterParams) port.LogsAdapter {
	if params == nil {
		params = &LogsFSAdapterParams{}
	}

	if params.ContainersPath == "" {
		params.ContainersPath = path.Join(storage.FSPath, "apps", "containers", "containers")
	}

	return &logsFSAdapter{
		loggers:   map[types.ContainerID]*ContainerLogger{},
		loggersMu: sync.RWMutex{},

		containersPath: params.ContainersPath,
	}
}

func (a *logsFSAdapter) Register(uuid types.ContainerID) error {
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

	a.loggersMu.Lock()
	defer a.loggersMu.Unlock()
	a.loggers[uuid] = &l

	err = l.Open()
	if err != nil {
		return err
	}

	return l.startCron()
}

func (a *logsFSAdapter) Unregister(uuid types.ContainerID) error {
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

	a.loggersMu.Lock()
	defer a.loggersMu.Unlock()
	delete(a.loggers, uuid)
	return nil
}

func (a *logsFSAdapter) Push(uuid types.ContainerID, line types.LogLine) {
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

func (a *logsFSAdapter) Pop(uuid types.ContainerID) (types.LogLine, error) {
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

func (a *logsFSAdapter) LoadBuffer(uuid types.ContainerID) ([]types.LogLine, error) {
	l, err := a.getLogger(uuid)
	if err != nil {
		return nil, err
	}
	return l.buffer, nil
}

func (a *logsFSAdapter) UnregisterAll() error {
	var ids []types.ContainerID

	a.loggersMu.RLock()
	for id := range a.loggers {
		ids = append(ids, id)
	}
	a.loggersMu.RUnlock()

	for _, id := range ids {
		log.Info("unregistering container logger", vlog.String("uuid", id.String()))
		err := a.Unregister(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *logsFSAdapter) getLogger(uuid types.ContainerID) (*ContainerLogger, error) {
	a.loggersMu.RLock()
	defer a.loggersMu.RUnlock()

	l, ok := a.loggers[uuid]
	if !ok {
		return nil, ErrLoggerNotFound
	}
	return l, nil
}

func (a *logsFSAdapter) dir(uuid types.ContainerID) string {
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
