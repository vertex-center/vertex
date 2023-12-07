package service

import (
	"errors"
	"io"
	"os"
	"path"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type dockerKernelService struct {
	dockerAdapter port.DockerAdapter
}

func NewDockerKernelService(dockerAdapter port.DockerAdapter) port.DockerService {
	return &dockerKernelService{
		dockerAdapter: dockerAdapter,
	}
}

func (s dockerKernelService) ListContainers() ([]types.DockerContainer, error) {
	return s.dockerAdapter.ListContainers()
}

func (s dockerKernelService) DeleteContainer(id string) error {
	return s.dockerAdapter.DeleteContainer(id)
}

func (s dockerKernelService) CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error) {
	return s.dockerAdapter.CreateContainer(options)
}

func (s dockerKernelService) StartContainer(id string) error {
	return s.dockerAdapter.StartContainer(id)
}

func (s dockerKernelService) StopContainer(id string) error {
	return s.dockerAdapter.StopContainer(id)
}

func (s dockerKernelService) InfoContainer(id string) (types.InfoContainerResponse, error) {
	return s.dockerAdapter.InfoContainer(id)
}

func (s dockerKernelService) LogsStdoutContainer(id string) (io.ReadCloser, error) {
	return s.dockerAdapter.LogsStdoutContainer(id)
}

func (s dockerKernelService) LogsStderrContainer(id string) (io.ReadCloser, error) {
	return s.dockerAdapter.LogsStderrContainer(id)
}

func (s dockerKernelService) WaitContainer(id string, cond types.WaitContainerCondition) error {
	return s.dockerAdapter.WaitContainer(id, cond)
}

func (s dockerKernelService) InfoImage(id string) (types.InfoImageResponse, error) {
	return s.dockerAdapter.InfoImage(id)
}

func (s dockerKernelService) PullImage(options types.PullImageOptions) (io.ReadCloser, error) {
	log.Info("pulling image", vlog.String("image", options.Image))
	return s.dockerAdapter.PullImage(options)
}

func (s dockerKernelService) BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error) {
	log.Info("building image", vlog.String("dockerfile", options.Dockerfile))
	return s.dockerAdapter.BuildImage(options)
}

func (s dockerKernelService) DeleteMounts(uuid string) error {
	volumesPath := path.Join("live_docker", "apps", "containers", "volumes", uuid)
	err := os.RemoveAll(volumesPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
