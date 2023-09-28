package adapter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type RunnerFSAdapter struct {
	commands map[uuid.UUID]*exec.Cmd
}

func NewRunnerFSAdapter() RunnerFSAdapter {
	return RunnerFSAdapter{
		commands: map[uuid.UUID]*exec.Cmd{},
	}
}

func (a RunnerFSAdapter) Delete(instance *types.Instance) error {
	return nil
}

func (a RunnerFSAdapter) Start(instance *types.Instance, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	service := instance.Service

	if a.commands[instance.UUID] != nil {
		log.Error(errors.New("runner already started"),
			vlog.String("name", service.Name),
		)
	}

	dir := a.getPath(*instance)
	executable := service.Methods.Script.Filename
	command := "./" + executable

	_, err = os.Stat(path.Join(dir, executable))
	if errors.Is(err, os.ErrNotExist) {
		err = fmt.Errorf("executable %s were not found", command)
		return
	} else if err != nil {
		return
	}

	a.commands[instance.UUID] = exec.Command(command)

	cmd := a.commands[instance.UUID]
	cmd.Dir = dir
	cmd.Env = os.Environ()

	var envFile *os.File
	envFile, err = os.Open(path.Join(a.getPath(*instance), ".env"))
	if err != nil {
		return
	}

	envScanner := bufio.NewScanner(envFile)
	for envScanner.Scan() {
		if envScanner.Text() == "" {
			continue
		}
		cmd.Env = append(cmd.Env, envScanner.Text())
	}

	cmd.Stdin = os.Stdin

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		return
	}

	setStatus(types.InstanceStatusRunning)

	err = cmd.Start()
	if err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		return
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Error(err,
				vlog.String("name", service.Name),
			)
		}
		setStatus(types.InstanceStatusOff)
	}()

	return
}

func (a RunnerFSAdapter) Stop(instance *types.Instance) error {
	cmd := a.commands[instance.UUID]

	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Force kill if the process continues

	delete(a.commands, instance.UUID)

	return nil
}

func (a RunnerFSAdapter) Info(instance types.Instance) (map[string]any, error) {
	return map[string]any{}, nil
}

func (a RunnerFSAdapter) CheckForUpdates(instance *types.Instance) error {
	//TODO implement me
	return nil
}

func (a RunnerFSAdapter) HasUpdateAvailable(instance types.Instance) (bool, error) {
	//TODO implement me
	return false, nil
}

func (a RunnerFSAdapter) getPath(instance types.Instance) string {
	return path.Join(storage.Path, "instances", instance.UUID.String())
}
