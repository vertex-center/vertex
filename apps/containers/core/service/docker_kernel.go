package service

import (
	"errors"
	"io"
	"os"
	"path"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vlog"
)

type dockerKernelService struct {
	adapter port.DockerAdapter
}

func NewDockerKernelService(adapter port.DockerAdapter) port.DockerService {
	return &dockerKernelService{adapter}
}

func (s dockerKernelService) ListContainers() ([]types.DockerContainer, error) {
	return s.adapter.ListContainers()
}

func (s dockerKernelService) DeleteContainer(id string) error {
	return s.adapter.DeleteContainer(id)
}

func (s dockerKernelService) CreateContainer(options types.CreateDockerContainerOptions) (types.CreateContainerResponse, error) {
	return s.adapter.CreateContainer(options)
}

func (s dockerKernelService) StartContainer(id string) error {
	return s.adapter.StartContainer(id)
}

func (s dockerKernelService) StopContainer(id string) error {
	return s.adapter.StopContainer(id)
}

func (s dockerKernelService) InfoContainer(id string) (types.InfoContainerResponse, error) {
	return s.adapter.InfoContainer(id)
}

func (s dockerKernelService) LogsStdoutContainer(id string) (io.ReadCloser, error) {
	return s.adapter.LogsStdoutContainer(id)
}

func (s dockerKernelService) LogsStderrContainer(id string) (io.ReadCloser, error) {
	return s.adapter.LogsStderrContainer(id)
}

func (s dockerKernelService) WaitContainer(id string, cond types.WaitContainerCondition) error {
	return s.adapter.WaitContainer(id, cond)
}

func (s dockerKernelService) InfoImage(id string) (types.InfoImageResponse, error) {
	return s.adapter.InfoImage(id)
}

func (s dockerKernelService) PullImage(options types.PullImageOptions) (io.ReadCloser, error) {
	log.Info("pulling image", vlog.String("image", options.Image))
	return s.adapter.PullImage(options)
}

func (s dockerKernelService) BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error) {
	log.Info("building image", vlog.String("dockerfile", options.Dockerfile))
	return s.adapter.BuildImage(options)
}

func (s dockerKernelService) DeleteMounts(uuid string) error {
	volumesPath := path.Join("live_docker", "apps", "containers", "volumes", uuid)
	err := os.RemoveAll(volumesPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (s dockerKernelService) CreateVolume(name string) (volume.Volume, error) {
	return s.adapter.CreateVolume(types.CreateVolumeOptions{
		Name: name,
	})
}

func (s dockerKernelService) DeleteVolume(name string) error {
	return s.adapter.DeleteVolume(name)
}
