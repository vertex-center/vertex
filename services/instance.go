package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
)

const (
	StatusOff     = "off"
	StatusRunning = "running"
	StatusError   = "error"
)

var (
	logger    = console.New("vertex::services-manager")
	instances = Instances{}
)

type Instance struct {
	Service

	uuid uuid.UUID
	cmd  *exec.Cmd
}

type Instances map[uuid.UUID]*Instance

func CreateFromDisk(uuid uuid.UUID) (*Instance, error) {
	data, err := os.ReadFile(path.Join("servers", uuid.String(), ".vertex", "service.json"))
	if err != nil {
		logger.Warn(fmt.Sprintf("service '%s' has no '.vertex/service.json' file", uuid))
	}

	var service Service
	err = json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}

	return &Instance{
		Service: service,
		uuid:    uuid,
	}, nil
}

func (s *Instance) Start() error {
	if s.cmd != nil {
		return errors.New("runner already started")
	}

	s.cmd = exec.Command(fmt.Sprintf("./%s", s.ID))
	s.cmd.Dir = path.Join("servers", s.uuid.String())

	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr
	s.cmd.Stdin = os.Stdin

	return s.cmd.Start()
}

func (s *Instance) Stop() error {
	err := s.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Remove runner from runners
	// TODO: Force kill if the process continues

	s.cmd = nil
	return nil
}

func (s *Instance) Status() string {
	if s.cmd == nil {
		return StatusOff
	}

	state := s.cmd.ProcessState
	if state == nil {
		return StatusRunning
	}

	return StatusError
}

func (s *Instance) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Instance
		Status string `json:"status"`
	}{
		Instance: *s,
		Status:   s.Status(),
	})
}

func Instantiate(uuid uuid.UUID) (*Instance, error) {
	if instances[uuid] != nil {
		return nil, fmt.Errorf("the service '%s' is already running", uuid)
	}

	instance, err := CreateFromDisk(uuid)
	if err != nil {
		return nil, err
	}

	instances[uuid] = instance

	return instance, nil
}

func ListInstances() Instances {
	return instances
}

func GetInstalled(uuid uuid.UUID) (*Instance, error) {
	instance := instances[uuid]
	if instance == nil {
		return nil, fmt.Errorf("the service '%s' is not instances", uuid)
	}
	return instance, nil
}

func IsInstantiated(uuid uuid.UUID) bool {
	return instances[uuid] != nil
}
