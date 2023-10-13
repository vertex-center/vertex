package services

import (
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type DockerKernelService struct {
	dockerAdapter types.DockerAdapterPort
}

func NewDockerKernelService(dockerAdapter types.DockerAdapterPort) DockerKernelService {
	return DockerKernelService{
		dockerAdapter: dockerAdapter,
	}
}

func (s DockerKernelService) ListContainers() ([]types.Container, error) {
	return s.dockerAdapter.ListContainers()
}

func (s DockerKernelService) DeleteContainer(id string) error {
	return s.dockerAdapter.DeleteContainer(id)
}

func (s DockerKernelService) CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error) {
	return s.dockerAdapter.CreateContainer(options)
}

func (s DockerKernelService) StartContainer(id string) error {
	return s.dockerAdapter.StartContainer(id)
}

func (s DockerKernelService) StopContainer(id string) error {
	return s.dockerAdapter.StopContainer(id)
}

func (s DockerKernelService) InfoContainer(id string) (types.InfoContainerResponse, error) {
	return s.dockerAdapter.InfoContainer(id)
}

func (s DockerKernelService) LogsStdoutContainer(id string) (io.ReadCloser, error) {
	return s.dockerAdapter.LogsStdoutContainer(id)
}

func (s DockerKernelService) LogsStderrContainer(id string) (io.ReadCloser, error) {
	return s.dockerAdapter.LogsStderrContainer(id)
}

func (s DockerKernelService) WaitContainer(id string, cond types.WaitContainerCondition) error {
	return s.dockerAdapter.WaitContainer(id, cond)
}

func (s DockerKernelService) InfoImage(id string) (types.InfoImageResponse, error) {
	return s.dockerAdapter.InfoImage(id)
}

func (s DockerKernelService) PullImage(options types.PullImageOptions) (io.ReadCloser, error) {
	log.Info("pulling image", vlog.String("image", options.Image))
	return s.dockerAdapter.PullImage(options)
}

func (s DockerKernelService) BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error) {
	log.Info("building image", vlog.String("dockerfile", options.Dockerfile))
	return s.dockerAdapter.BuildImage(options)
}
