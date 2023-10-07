package adapter

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type DockerCliAdapter struct {
	cli *client.Client
}

func NewDockerCliAdapter() DockerCliAdapter {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Warn("couldn't connect with the Docker cli.",
			vlog.String("error", err.Error()),
		)

		return DockerCliAdapter{}
	}

	return DockerCliAdapter{
		cli: cli,
	}
}

func (a DockerCliAdapter) ListContainers() ([]types.Container, error) {
	res, err := a.cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var containers []types.Container
	for _, c := range res {
		containers = append(containers, types.NewContainer(c))
	}
	return containers, nil
}

func (a DockerCliAdapter) DeleteContainer(id string) error {
	return a.cli.ContainerRemove(context.Background(), id, dockertypes.ContainerRemoveOptions{})
}

func (a DockerCliAdapter) CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error) {
	config := container.Config{
		Image:        options.ImageName,
		ExposedPorts: options.ExposedPorts,
		Env:          options.Env,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          options.Cmd,
	}

	hostConfig := container.HostConfig{
		Binds:        options.Binds,
		PortBindings: options.PortBindings,
		CapAdd:       options.CapAdd,
		Sysctls:      options.Sysctls,
	}

	res, err := a.cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, options.ContainerName)
	if err != nil {
		return types.CreateContainerResponse{}, err
	}

	return types.CreateContainerResponse{
		ID:       res.ID,
		Warnings: res.Warnings,
	}, nil
}

func (a DockerCliAdapter) StartContainer(id string) error {
	return a.cli.ContainerStart(context.Background(), id, dockertypes.ContainerStartOptions{})
}

func (a DockerCliAdapter) StopContainer(id string) error {
	return a.cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (a DockerCliAdapter) InfoContainer(id string) (types.InfoContainerResponse, error) {
	info, err := a.cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return types.InfoContainerResponse{}, err
	}
	return types.InfoContainerResponse{
		ID:       info.ID,
		Name:     info.Name,
		Platform: info.Platform,
		Image:    info.Image,
	}, nil
}

func (a DockerCliAdapter) LogsStdoutContainer(id string) (io.ReadCloser, error) {
	return a.cli.ContainerLogs(context.Background(), id, dockertypes.ContainerLogsOptions{
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "0",
	})
}

func (a DockerCliAdapter) LogsStderrContainer(id string) (io.ReadCloser, error) {
	return a.cli.ContainerLogs(context.Background(), id, dockertypes.ContainerLogsOptions{
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "0",
	})
}

func (a DockerCliAdapter) WaitContainer(id string, cond types.WaitContainerCondition) error {
	statusCh, errCh := a.cli.ContainerWait(context.Background(), id, container.WaitCondition(cond))

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}

func (a DockerCliAdapter) InfoImage(id string) (types.InfoImageResponse, error) {
	info, _, err := a.cli.ImageInspectWithRaw(context.Background(), id)
	if err != nil {
		return types.InfoImageResponse{}, nil
	}
	return types.InfoImageResponse{
		ID:           info.ID,
		Architecture: info.Architecture,
		OS:           info.Os,
		Size:         info.Size,
		Tags:         info.RepoTags,
	}, nil
}

func (a DockerCliAdapter) PullImage(options types.PullImageOptions) (io.ReadCloser, error) {
	return a.cli.ImagePull(context.Background(), options.Image, dockertypes.ImagePullOptions{})
}

func (a DockerCliAdapter) BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error) {
	buildOptions := dockertypes.ImageBuildOptions{
		Dockerfile: options.Dockerfile,
		Tags:       []string{options.Name},
		Remove:     true,
	}

	reader, err := archive.TarWithOptions(options.Dir, &archive.TarOptions{
		ExcludePatterns: []string{".git/**/*"},
	})
	if err != nil {
		return dockertypes.ImageBuildResponse{}, err
	}

	return a.cli.ImageBuild(context.Background(), reader, buildOptions)
}
