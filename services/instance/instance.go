package instance

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/services"
)

var logger = console.New("vertex::instance")

const (
	StatusOff     = "off"
	StatusRunning = "running"
	StatusError   = "error"
)

const (
	EventStdout       = "stdout"
	EventStderr       = "stderr"
	EventStatusChange = "status_change"
)

type Event struct {
	Name string
	Data string
}

type Instance struct {
	services.Service

	Status string `json:"status"`
	Logs   Logs   `json:"logs"`

	uuid uuid.UUID
	cmd  *exec.Cmd

	listeners map[uuid.UUID]chan Event
}

func (i *Instance) Start() error {
	if i.cmd != nil {
		logger.Error(fmt.Errorf("runner %s already started", i.Name))
	}

	i.cmd = exec.Command(fmt.Sprintf("./%s", i.ID))
	i.cmd.Dir = path.Join("servers", i.uuid.String())

	i.cmd.Stdin = os.Stdin

	stdoutReader, err := i.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := i.cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			line := i.Logs.Add(LogLine{
				Kind:    LogKindOut,
				Message: stdoutScanner.Text(),
			})

			data, err := json.Marshal(line)
			if err != nil {
				logger.Error(err)
			}

			for _, listener := range i.listeners {
				listener <- Event{
					Name: EventStdout,
					Data: string(data),
				}
			}
		}
	}()

	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			line := i.Logs.Add(LogLine{
				Kind:    LogKindErr,
				Message: stderrScanner.Text(),
			})

			data, err := json.Marshal(line)
			if err != nil {
				logger.Error(err)
			}

			for _, listener := range i.listeners {
				listener <- Event{
					Name: EventStderr,
					Data: string(data),
				}
			}
		}
	}()

	i.setStatus(StatusRunning)

	err = i.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := i.cmd.Wait()
		if err != nil {
			logger.Error(fmt.Errorf("%s: %v", i.Service.Name, err))
		}
		i.setStatus(StatusOff)
	}()

	return nil
}

func (i *Instance) Stop() error {
	err := i.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Force kill if the process continues

	i.cmd = nil

	return nil
}

func (i *Instance) setStatus(status string) {
	i.Status = status

	for _, listener := range i.listeners {
		listener <- Event{
			Name: EventStatusChange,
			Data: status,
		}
	}
}

func (i *Instance) Register(channel chan Event) uuid.UUID {
	id := uuid.New()
	i.listeners[id] = channel
	logger.Log(fmt.Sprintf("channel %s registered to instance uuid=%s", id, i.uuid))
	return id
}

func (i *Instance) Unregister(uuid uuid.UUID) {
	delete(i.listeners, uuid)
	logger.Log(fmt.Sprintf("channel %s unregistered from instance uuid=%s", uuid, i.uuid))
}

func CreateFromDisk(instanceUUID uuid.UUID) (*Instance, error) {
	data, err := os.ReadFile(path.Join("servers", instanceUUID.String(), ".vertex", "service.json"))
	if err != nil {
		logger.Warn(fmt.Sprintf("service '%s' has no '.vertex/service.json' file", instanceUUID))
	}

	var service services.Service
	err = json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}

	return &Instance{
		Service:   service,
		Status:    StatusOff,
		Logs:      Logs{},
		uuid:      instanceUUID,
		listeners: map[uuid.UUID]chan Event{},
	}, nil
}
