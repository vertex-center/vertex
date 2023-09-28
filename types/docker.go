package types

import (
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type Container struct {
	ID      string   `json:"id,omitempty"`
	ImageID string   `json:"image_id,omitempty"`
	Names   []string `json:"names,omitempty"`
	Mounts  []Mount  `json:"mounts,omitempty"`
}

type Mount struct {
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
}

type CreateContainerOptions struct {
	ImageName     string            `json:"image_name,omitempty"`
	ContainerName string            `json:"container_name,omitempty"`
	ExposedPorts  nat.PortSet       `json:"exposed_ports,omitempty"`
	PortBindings  nat.PortMap       `json:"port_bindings,omitempty"`
	Binds         []string          `json:"binds,omitempty"`
	Env           []string          `json:"env,omitempty"`
	CapAdd        []string          `json:"cap_add,omitempty"`
	Sysctls       map[string]string `json:"sysctls,omitempty"`
}

type BuildImageOptions struct {
	Dir        string `json:"dir,omitempty"`
	Name       string `json:"name,omitempty"`
	Dockerfile string `json:"dockerfile,omitempty"`
}

type PullImageOptions struct {
	Image string `json:"image,omitempty"`
}

type CreateContainerResponse struct {
	ID       string   `json:"id,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type InfoContainerResponse struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Platform string `json:"platform,omitempty"`
	Image    string `json:"image,omitempty"`
}

type InfoImageResponse struct {
	ID           string   `json:"id,omitempty"`
	Architecture string   `json:"architecture,omitempty"`
	OS           string   `json:"os,omitempty"`
	Size         int64    `json:"size,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

type WaitContainerCondition container.WaitCondition

type DockerAdapterPort interface {
	ListContainers() ([]Container, error)
	DeleteContainer(id string) error
	CreateContainer(options CreateContainerOptions) (CreateContainerResponse, error)
	StartContainer(id string) error
	StopContainer(id string) error
	InfoContainer(id string) (InfoContainerResponse, error)
	LogsStdoutContainer(id string) (io.ReadCloser, error)
	LogsStderrContainer(id string) (io.ReadCloser, error)
	WaitContainer(id string, cond WaitContainerCondition) error
	InfoImage(id string) (InfoImageResponse, error)
	PullImage(options PullImageOptions) (io.ReadCloser, error)
	BuildImage(options BuildImageOptions) (dockertypes.ImageBuildResponse, error)
}

func NewContainer(c dockertypes.Container) Container {
	return Container{
		ID:      c.ID,
		ImageID: c.ImageID,
		Names:   c.Names,
		Mounts:  NewMounts(c.Mounts),
	}
}

func NewMounts(m []dockertypes.MountPoint) []Mount {
	mounts := make([]Mount, len(m))
	for i, v := range m {
		mounts[i] = NewMount(v)
	}
	return mounts
}

func NewMount(m dockertypes.MountPoint) Mount {
	return Mount{
		Source:      m.Source,
		Destination: m.Destination,
	}
}
