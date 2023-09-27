package services

import (
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/types"
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

func (s DockerKernelService) LogsContainer(id string) (io.ReadCloser, error) {
	return s.dockerAdapter.LogsContainer(id)
}

func (s DockerKernelService) WaitContainer(id string, cond types.WaitContainerCondition) error {
	return s.dockerAdapter.WaitContainer(id, cond)
}

func (s DockerKernelService) InfoImage(id string) (types.InfoImageResponse, error) {
	return s.dockerAdapter.InfoImage(id)
}

func (s DockerKernelService) PullImage(options types.PullImageOptions) (io.ReadCloser, error) {
	return s.dockerAdapter.PullImage(options)
}

func (s DockerKernelService) BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error) {
	return s.dockerAdapter.BuildImage(options)
}
