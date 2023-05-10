package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
)

type FSRunnerRepository struct{}

func NewFSRunnerRepository() FSRunnerRepository {
	return FSRunnerRepository{}
}

func (r FSRunnerRepository) Delete(instance *types.Instance) error {
	return nil
}

func (r FSRunnerRepository) Start(instance *types.Instance) error {
	if instance.Cmd != nil {
		logger.Error(errors.New("runner already started")).
			AddKeyValue("name", instance.Name).
			Print()
	}

	dir := r.getPath(*instance)
	executable := instance.ID
	command := "./" + instance.ID

	// Try to find the executable
	// For a service of ID=vertex-id, the executable can be:
	// - vertex-id
	// - script-filename.sh
	_, err := os.Stat(path.Join(dir, executable))
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Stat(path.Join(dir, instance.Methods.Script.Filename))
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("executables %s and %s were not found", instance.ID, instance.Methods.Script.Filename)
		} else if err != nil {
			return err
		}
		command = fmt.Sprintf("./%s", instance.Methods.Script.Filename)
	} else if err != nil {
		return err
	}

	instance.Cmd = exec.Command(command)
	instance.Cmd.Dir = dir

	instance.Cmd.Stdin = os.Stdin

	stdoutReader, err := instance.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := instance.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			instance.PushLogLine(&types.LogLine{
				Kind:    types.LogKindOut,
				Message: stdoutScanner.Text(),
			})
		}
	}()

	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			instance.PushLogLine(&types.LogLine{
				Kind:    types.LogKindErr,
				Message: stderrScanner.Text(),
			})
		}
	}()

	instance.SetStatus(types.InstanceStatusRunning)

	err = instance.Cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := instance.Cmd.Wait()
		if err != nil {
			logger.Error(err).
				AddKeyValue("name", instance.Service.Name).
				Print()
		}
		instance.SetStatus(types.InstanceStatusOff)
	}()

	return nil
}

func (r FSRunnerRepository) Stop(instance *types.Instance) error {
	err := instance.Cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	// TODO: Force kill if the process continues

	instance.Cmd = nil

	return nil
}

func (r FSRunnerRepository) Info(instance types.Instance) (map[string]any, error) {
	return map[string]any{}, nil
}

func (r FSRunnerRepository) getPath(instance types.Instance) string {
	return path.Join(storage.PathInstances, instance.UUID.String())
}
