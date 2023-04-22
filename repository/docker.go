package repository

import (
	"bufio"
	"context"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type DockerRepository struct {
	cli *client.Client
}

func NewDockerRepository() DockerRepository {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Warn("Couldn't connect Docker cli.")
		return DockerRepository{}
	}

	return DockerRepository{
		cli: cli,
	}
}

func (r DockerRepository) RemoveContainer(id string) error {
	return r.cli.ContainerRemove(context.Background(), id, dockertypes.ContainerRemoveOptions{})
}

func (r DockerRepository) BuildImage(instancePath string, imageName string) error {
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
		logger.Log(scanner.Text())
	}

	logger.Log("Docker build: success.")
	return nil
}

func (r DockerRepository) CreateContainer(imageName string, containerName string) (string, error) {
	res, err := r.cli.ContainerCreate(context.Background(), &container.Config{
		Image: imageName,
	}, nil, nil, nil, containerName)
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