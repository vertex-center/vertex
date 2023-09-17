package adapter

import (
	"bufio"
	"errors"
	"fmt"
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

func (a RunnerFSAdapter) Start(instance *types.Instance, onLog func(msg string), onErr func(msg string), setStatus func(status string)) error {
	service := instance.Service

	if a.commands[instance.UUID] != nil {
		log.Error(errors.New("runner already started"),
			vlog.String("name", service.Name),
		)
	}

	dir := a.getPath(*instance)
	executable := service.Methods.Script.Filename
	command := "./" + executable

	_, err := os.Stat(path.Join(dir, executable))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("executable %s were not found", command)
	} else if err != nil {
		return err
	}

	a.commands[instance.UUID] = exec.Command(command)

	cmd := a.commands[instance.UUID]
	cmd.Dir = dir
	cmd.Env = os.Environ()

	envFile, err := os.Open(path.Join(a.getPath(*instance), ".env"))
	if err != nil {
		return err
	}

	envScanner := bufio.NewScanner(envFile)
	for envScanner.Scan() {
		if envScanner.Text() == "" {
			continue
		}
		cmd.Env = append(cmd.Env, envScanner.Text())
	}

	cmd.Stdin = os.Stdin

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			onLog(stdoutScanner.Text())
		}
	}()

	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			onErr(stderrScanner.Text())
		}
	}()

	setStatus(types.InstanceStatusRunning)

	err = cmd.Start()
	if err != nil {
		return err
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

	return nil
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
