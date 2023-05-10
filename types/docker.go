package types

import "github.com/docker/go-connections/nat"

type DockerContainerInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Platform string `json:"platform"`
}

// FIXME: The DockerRepository interface should not handle that much

type DockerRepository interface {
	RemoveContainer(id string) error
	BuildImageFromName(imageName string, onMsg func(msg string)) error
	BuildImageFromDockerfile(instancePath string, imageName string, onMsg func(msg string)) error
	CreateContainer(imageName string, containerName string, exposedPorts nat.PortSet, portBindings nat.PortMap, binds []string) (string, error)
	StartContainer(id string) error
	StopContainer(id string) error
	GetContainerID(containerName string) (string, error)
	GetContainerInfo(containerName string) (*DockerContainerInfo, error)
}
