package runners

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

const (
	StatusOff     = "off"
	StatusRunning = "running"
	StatusError   = "error"
)

var runners = map[string]*Runner{}

type Runner struct {
	ServiceID string
	cmd       *exec.Cmd
}

func NewRunner(serviceID string) (*Runner, error) {
	if runners[serviceID] != nil {
		return nil, fmt.Errorf("the service '%s' is already running", serviceID)
	}

	runner := &Runner{
		ServiceID: serviceID,
	}
	runners[serviceID] = runner
	return runner, nil
}

func (r *Runner) Start() error {
	if r.cmd != nil {
		return errors.New("runner already started")
	}

	r.cmd = exec.Command(fmt.Sprintf("./%s", r.ServiceID))
	r.cmd.Dir = path.Join("servers", r.ServiceID)

	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr
	r.cmd.Stdin = os.Stdin

	return r.cmd.Start()
}

func (r *Runner) Stop() error {
	err := r.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Remove runner from runners
	// TODO: Force kill if the process continues

	r.cmd = nil
	return nil
}

func (r *Runner) Status() string {
	if r.cmd == nil {
		return StatusOff
	}

	state := r.cmd.ProcessState
	if state == nil {
		return StatusRunning
	}

	return StatusError
}

func GetRunner(serviceID string) (*Runner, error) {
	runner := runners[serviceID]
	if runner == nil {
		return nil, fmt.Errorf("the runner '%s' was nos found", serviceID)
	}
	return runner, nil
}
