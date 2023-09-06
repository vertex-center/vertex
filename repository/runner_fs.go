package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
)

type RunnerFSRepository struct {
	commands map[uuid.UUID]*exec.Cmd
}

func NewRunnerFSRepository() RunnerFSRepository {
	return RunnerFSRepository{
		commands: map[uuid.UUID]*exec.Cmd{},
	}
}

func (r RunnerFSRepository) Delete(instance *types.Instance) error {
	return nil
}

func (r RunnerFSRepository) Start(instance *types.Instance, onLog func(msg string), onErr func(msg string), setStatus func(status string)) error {
	if r.commands[instance.UUID] != nil {
		logger.Error(errors.New("runner already started")).
			AddKeyValue("name", instance.Name).
			Print()
	}

	dir := r.getPath(*instance)
	executable := instance.Methods.Script.Filename
	command := "./" + executable

	_, err := os.Stat(path.Join(dir, executable))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("executable %s were not found", command)
	} else if err != nil {
		return err
	}

	r.commands[instance.UUID] = exec.Command(command)

	cmd := r.commands[instance.UUID]
	cmd.Dir = dir
	cmd.Env = os.Environ()

	envFile, err := os.Open(path.Join(r.getPath(*instance), ".env"))
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
			logger.Error(err).
				AddKeyValue("name", instance.Service.Name).
				Print()
		}
		setStatus(types.InstanceStatusOff)
	}()

	return nil
}

func (r RunnerFSRepository) Stop(instance *types.Instance) error {
	cmd := r.commands[instance.UUID]

	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Force kill if the process continues

	delete(r.commands, instance.UUID)

	return nil
}

func (r RunnerFSRepository) Info(instance types.Instance) (map[string]any, error) {
	return map[string]any{}, nil
}

func (r RunnerFSRepository) CheckForUpdates(instance *types.Instance) error {
	//TODO implement me
	return nil
}

func (r RunnerFSRepository) HasUpdateAvailable(instance types.Instance) (bool, error) {
	//TODO implement me
	return false, nil
}

func (r RunnerFSRepository) getPath(instance types.Instance) string {
	return path.Join(storage.Path, "instances", instance.UUID.String())
}
