package instance

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/storage"
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

type Metadata struct {
	UseDocker   bool `json:"use_docker"`
	UseReleases bool `json:"use_releases"`
}

var (
	errContainerNotFound = errors.New("container not found")
)

type Instance struct {
	services.Service
	Metadata

	Status       string       `json:"status"`
	Logs         Logs         `json:"logs"`
	EnvVariables EnvVariables `json:"env"`

	UUID uuid.UUID `json:"uuid"`
	cmd  *exec.Cmd

	listeners map[uuid.UUID]chan Event
}

func (i *Instance) dockerImageName() string {
	return "vertex_image_" + i.UUID.String()
}

func (i *Instance) dockerContainerName() string {
	return "VERTEX_CONTAINER_" + i.UUID.String()
}

func (i *Instance) dockerContainerID(cli *client.Client) (string, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}

	var containerID string

	for _, c := range containers {
		name := c.Names[0]
		if name == "/"+i.dockerContainerName() {
			containerID = c.ID
			break
		}
	}

	if containerID == "" {
		return "", errContainerNotFound
	}

	return containerID, nil
}

func (i *Instance) Start() error {
	var err error
	if i.UseDocker {
		err = i.startWithDocker()
	} else {
		err = i.startManually()
	}
	return err
}

func (i *Instance) startManually() error {
	if i.cmd != nil {
		logger.Error(fmt.Errorf("runner %s already started", i.Name))
	}

	dir := path.Join(storage.PathInstances, i.UUID.String())
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

func (i *Instance) startWithDocker() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	imageName := i.dockerImageName()
	containerName := i.dockerContainerName()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	}

	reader, err := archive.TarWithOptions(path.Join(storage.PathInstances, i.UUID.String()), &archive.TarOptions{
		ExcludePatterns: []string{".git/**/*"},
	})
	if err != nil {
		return err
	}

	i.setStatus(StatusRunning)

	res, err := cli.ImageBuild(context.Background(), reader, buildOptions)
	if err != nil {
		i.setStatus(StatusOff)
		return err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		if scanner.Err() != nil {
			i.setStatus(StatusOff)
			return scanner.Err()
		}
		logger.Log(scanner.Text())
	}

	logger.Log("Docker build: success.")

	id, err := i.dockerContainerID(cli)
	if err == errContainerNotFound {
		logger.Log(fmt.Sprintf("container %s doesn't exists, create it.", containerName))
		res, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: imageName,
		}, nil, nil, nil, containerName)
		if err != nil {
			i.setStatus(StatusOff)
			return err
		}
		id = res.ID
	} else if err != nil {
		i.setStatus(StatusOff)
		return err
	}

	logger.Log("starting container...")

	err = cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		i.setStatus(StatusOff)
		return err
	}
	return nil
}

func (i *Instance) Stop() error {
	var err error
	if i.UseDocker {
		err = i.stopWithDocker()
	} else {
		err = i.stopManually()
	}
	i.setStatus(StatusOff)
	return err
}

func (i *Instance) stopWithDocker() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	id, err := i.dockerContainerID(cli)
	if err != nil {
		return err
	}

	return cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (i *Instance) stopManually() error {
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
	logger.Log(fmt.Sprintf("channel %s registered to instance uuid=%s", id, i.UUID))
	return id
}

func (i *Instance) Unregister(uuid uuid.UUID) {
	delete(i.listeners, uuid)
	logger.Log(fmt.Sprintf("channel %s unregistered from instance uuid=%s", uuid, i.UUID))
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

	err := os.RemoveAll(path.Join(storage.PathInstances, i.UUID.String()))
	if err != nil {
		return fmt.Errorf("failed to delete server uuid=%s: %v", i.UUID, err)
	}
	return nil
}

func CreateFromDisk(instanceUUID uuid.UUID) (*Instance, error) {
	service, err := services.ReadFromDisk(path.Join(storage.PathInstances, instanceUUID.String()))
	if err != nil {
		return nil, err
	}

	meta := Metadata{
		UseDocker:   false,
		UseReleases: false,
	}

	metaPath := path.Join(storage.PathInstances, instanceUUID.String(), ".vertex", "instance_metadata.json")
	metaBytes, err := os.ReadFile(metaPath)

	if errors.Is(err, os.ErrNotExist) {
		logger.Log("instance_metadata.json not found. using default.")
	} else if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(metaBytes, &meta)
		if err != nil {
			return nil, err
		}
	}

	i := &Instance{
		Service:      *service,
		Metadata:     meta,
		Status:       StatusOff,
		Logs:         Logs{},
		EnvVariables: *NewEnvVariables(),
		UUID:         instanceUUID,
		listeners:    map[uuid.UUID]chan Event{},
	}

	err = i.LoadEnvFromDisk()
	return i, err
}

func (i *Instance) WriteMetadata() error {
	metaPath := path.Join(storage.PathInstances, i.UUID.String(), ".vertex", "instance_metadata.json")

	metaBytes, err := json.MarshalIndent(i.Metadata, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(metaPath, metaBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
