package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/vertex-center/vertex-core-golang/console"
)

const (
	StatusOff     = "off"
	StatusRunning = "running"
	StatusError   = "error"
)

var (
	logger    = console.New("vertex::services-manager")
	installed = map[string]*InstalledService{}
)

type InstalledService struct {
	Service
	cmd *exec.Cmd
}

func FromDisk(serviceID string) (*InstalledService, error) {
	data, err := os.ReadFile(path.Join("servers", serviceID, ".vertex", "service.json"))
	if err != nil {
		logger.Warn(fmt.Sprintf("service '%s' has no '.vertex/service.json' file", serviceID))
	}

	var service Service
	err = json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}

	return &InstalledService{
		Service: service,
	}, nil
}

func (s *InstalledService) Start() error {
	if s.cmd != nil {
		return errors.New("runner already started")
	}

	s.cmd = exec.Command(fmt.Sprintf("./%s", s.Service.ID))
	s.cmd.Dir = path.Join("servers", s.Service.ID)

	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr
	s.cmd.Stdin = os.Stdin

	return s.cmd.Start()
}

func (s *InstalledService) Stop() error {
	err := s.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Remove runner from runners
	// TODO: Force kill if the process continues

	s.cmd = nil
	return nil
}

func (s *InstalledService) Status() string {
	if s.cmd == nil {
		return StatusOff
	}

	state := s.cmd.ProcessState
	if state == nil {
		return StatusRunning
	}

	return StatusError
}

func (s *InstalledService) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InstalledService
		Status string `json:"status"`
	}{
		InstalledService: *s,
		Status:           s.Status(),
	})
}

func (s Service) Instantiate() (*InstalledService, error) {
	if installed[s.ID] != nil {
		return nil, fmt.Errorf("the service '%s' is already running", s.ID)
	}

	is, err := FromDisk(s.ID)
	if err != nil {
		return nil, err
	}

	installed[s.ID] = is

	return is, nil
}

func ListInstalled() map[string]*InstalledService {
	return installed
}

func GetInstalled(serviceID string) (*InstalledService, error) {
	service := installed[serviceID]
	if service == nil {
		return nil, fmt.Errorf("the service '%s' is not installed", serviceID)
	}
	return service, nil
}

func (s Service) IsInstalled() bool {
	return installed[s.ID] != nil
}
