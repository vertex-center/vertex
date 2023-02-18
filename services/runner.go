package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

var runners []*Runner

type Runner struct {
	Service Service
	cmd     *exec.Cmd
}

func NewRunner(service Service) *Runner {
	runner := &Runner{
		Service: service,
	}
	runners = append(runners, runner)
	return runner
}

func (r *Runner) Start() error {
	if r.cmd != nil {
		return errors.New("runner already started")
	}

	r.cmd = exec.Command(fmt.Sprintf("./%s", r.Service.ID))
	r.cmd.Dir = path.Join("servers", r.Service.ID)

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

func GetRunner() *Runner {
	// TODO
	return runners[0]
}
