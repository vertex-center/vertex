package repository

import (
	"bufio"
	"context"
	"encoding/json"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/vertex-center/vertex/logger"
	"github.com/vertex-center/vertex/types"
)

type DockerRepository struct {
	cli *client.Client
}

type dockerMessage struct {
	Stream string `json:"stream"`
}

func NewDockerRepository() DockerRepository {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Warn("couldn't connect with the Docker cli.").
			AddKeyValue("error", err.Error()).
			Print()

		return DockerRepository{}
	}

	return DockerRepository{
		cli: cli,
	}
}

func (r DockerRepository) RemoveContainer(id string) error {
	return r.cli.ContainerRemove(context.Background(), id, dockertypes.ContainerRemoveOptions{})
}

func (r DockerRepository) BuildImageFromName(imageName string, onMsg func(msg string)) error {
	res, err := r.cli.ImagePull(context.Background(), imageName, dockertypes.ImagePullOptions{})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(res)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		onMsg(scanner.Text())
	}

	return nil
}

func (r DockerRepository) BuildImageFromDockerfile(instancePath string, imageName string, onMsg func(msg string)) error {
	buildOptions := dockertypes.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	}

	reader, err := archive.TarWithOptions(instancePath, &archive.TarOptions{
		ExcludePatterns: []string{".git/**/*"},
	})
	if err != nil {
		return err
	}

	res, err := r.cli.ImageBuild(context.Background(), reader, buildOptions)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		msg := dockerMessage{}
		err := json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			logger.Warn("Failed to parse message:").
				AddKeyValue("msg", scanner.Text()).
				Print()
		} else {
			if msg.Stream != "" {
				onMsg(msg.Stream)
			}
		}
	}

	logger.Log("Docker build: success.").Print()
	return nil
}

func (r DockerRepository) CreateContainer(imageName string, containerName string, exposedPorts nat.PortSet, portBindings nat.PortMap, binds []string) (string, error) {
	config := container.Config{
		Image:        imageName,
		ExposedPorts: exposedPorts,
	}

	hostConfig := container.HostConfig{
		Binds:        binds,
		PortBindings: portBindings,
	}

	res, err := r.cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, containerName)
	for _, warn := range res.Warnings {
		logger.Warn("warning while creating container").
			AddKeyValue("warning", warn).
			Print()
	}
	return res.ID, err
}

func (r DockerRepository) StartContainer(id string) error {
	return r.cli.ContainerStart(context.Background(), id, dockertypes.ContainerStartOptions{})
}

func (r DockerRepository) StopContainer(id string) error {
	return r.cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (r DockerRepository) GetContainerID(containerName string) (string, error) {
	containers, err := r.cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}

	var containerID string

	for _, c := range containers {
		name := c.Names[0]
		if name == "/"+containerName {
			containerID = c.ID
			break
		}
	}

	if containerID == "" {
		return "", ErrContainerNotFound
	}

	return containerID, nil
}

func (r DockerRepository) GetContainerInfo(containerName string) (*types.DockerContainerInfo, error) {
	id, err := r.GetContainerID(containerName)
	if err != nil {
		return nil, err
	}
	info, err := r.cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &types.DockerContainerInfo{
		ID:       info.ID,
		Name:     info.Name,
		Image:    info.Image,
		Platform: info.Platform,
	}, nil
}
