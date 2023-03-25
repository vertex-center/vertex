package instance

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

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

type EnvVariables = map[string]string

type Instance struct {
	services.Service

	Status       string       `json:"status"`
	Logs         Logs         `json:"logs"`
	EnvVariables EnvVariables `json:"env"`

	uuid uuid.UUID
	cmd  *exec.Cmd

	listeners map[uuid.UUID]chan Event
}

func (i *Instance) Start() error {
	if i.cmd != nil {
		logger.Error(fmt.Errorf("runner %s already started", i.Name))
	}

	dir := path.Join("servers", i.uuid.String())
	executable := i.ID
	command := "./" + i.ID

	// Try to find the executable
	// For a service of ID=vertex-id, the executable can be:
	// - vertex-id
	// - vertex-id.sh
	_, err := os.Stat(path.Join(dir, executable))
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Stat(path.Join(dir, executable+".sh"))
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("the executable %s (or %s.sh) was not found at path", i.ID, i.ID)
		} else if err != nil {
			return err
		}
		command = fmt.Sprintf("./%s.sh", i.ID)
	} else if err != nil {
		return err
	}

	i.cmd = exec.Command(command)
	i.cmd.Dir = dir

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
	err := i.cmd.Process.Signal(os.Kill)
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

func (i *Instance) IsRunning() bool {
	return i.Status == StatusRunning
}

func (i *Instance) Delete() error {
	if i.IsRunning() {
		err := i.Stop()
		if err != nil {
			return err
		}
	}

	err := os.RemoveAll(path.Join("servers", i.uuid.String()))
	if err != nil {
		return fmt.Errorf("failed to delete server uuid=%s: %v", i.uuid, err)
	}
	return nil
}

func CreateFromDisk(instanceUUID uuid.UUID) (*Instance, error) {
	service, err := services.ReadFromDisk(path.Join("servers", instanceUUID.String()))
	if err != nil {
		return nil, err
	}

	env, err := readEnv(path.Join("servers", instanceUUID.String(), ".env"))
	if err != nil {
		return nil, err
	}

	return &Instance{
		Service:      *service,
		Status:       StatusOff,
		Logs:         Logs{},
		EnvVariables: env,
		uuid:         instanceUUID,
		listeners:    map[uuid.UUID]chan Event{},
	}, nil
}

func readEnv(filepath string) (EnvVariables, error) {
	var variables = EnvVariables{}

	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return EnvVariables{}, nil
	}
	if err != nil {
		return EnvVariables{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) < 2 {
			return EnvVariables{}, errors.New("failed to read .env")
		}

		variables[line[0]] = line[1]
	}

	return variables, nil
}
