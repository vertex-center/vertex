package types

import (
	"errors"

	"github.com/google/uuid"
)

const (
	ContainerStatusOff      = "off"
	ContainerStatusBuilding = "building"
	ContainerStatusStarting = "starting"
	ContainerStatusRunning  = "running"
	ContainerStatusStopping = "stopping"
	ContainerStatusError    = "error"
)

const (
	ContainerInstallMethodDocker = "docker"
)

var (
	ErrContainerNotFound     = errors.New("container not found")
	ErrContainerStillRunning = errors.New("container still running")
)

type Container struct {
	ContainerSettings

	Service Service               `json:"service"`
	UUID    uuid.UUID             `json:"uuid"`
	Status  string                `json:"status"`
	Env     ContainerEnvVariables `json:"environment,omitempty"`

	Update        *ContainerUpdate `json:"update,omitempty"`
	ServiceUpdate ServiceUpdate    `json:"service_update,omitempty"`

	CacheVersions []string `json:"cache_versions,omitempty"`
}

type ContainerSearchQuery struct {
	Tags     *[]string `json:"tags,omitempty"`
	Features *[]string `json:"features,omitempty"`
}

type ContainerUpdate struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

type DownloadProgress struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Current int64  `json:"current,omitempty"`
	Total   int64  `json:"total,omitempty"`
}

func NewContainer(id uuid.UUID, service Service) Container {
	return Container{
		Service: service,
		UUID:    id,
		Status:  ContainerStatusOff,
		Env:     map[string]string{},
	}
}

type ContainerAdapterPort interface {
	Create(uuid uuid.UUID) error
	Delete(uuid uuid.UUID) error

	GetAll() ([]uuid.UUID, error)
}

func (i *Container) DockerImageVertexName() string {
	return "vertex_image_" + i.UUID.String()
}

func (i *Container) DockerContainerName() string {
	return "VERTEX_CONTAINER_" + i.UUID.String()
}

func (i *Container) IsRunning() bool {
	return i.Status != ContainerStatusOff && i.Status != ContainerStatusError
}

func (i *Container) IsBusy() bool {
	return i.Status == ContainerStatusBuilding || i.Status == ContainerStatusStarting || i.Status == ContainerStatusStopping
}

func (i *Container) IsDockerized() bool {
	return i.ContainerSettings.InstallMethod != nil && *i.ContainerSettings.InstallMethod == ContainerInstallMethodDocker
}

func (i *Container) LaunchOnStartup() bool {
	launchOnStartup := i.ContainerSettings.LaunchOnStartup
	if launchOnStartup != nil && !*launchOnStartup {
		return false
	}
	return true
}

func (i *Container) ResetDefaultEnv() {
	i.Env = ContainerEnvVariables{}
	for _, env := range i.Service.Env {
		i.Env[env.Name] = env.Default
	}
}

func (i *Container) HasFeature(featureType string) bool {
	if i.Service.Features == nil {
		return false
	}

	if i.Service.Features.Databases != nil {
		for _, db := range *i.Service.Features.Databases {
			if db.Type == featureType {
				return true
			}
		}
	}

	return false
}

func (i *Container) HasFeatureIn(featureTypes []string) bool {
	if featureTypes == nil {
		return true
	}

	for _, featureType := range featureTypes {
		if i.HasFeature(featureType) {
			return true
		}
	}

	return false
}

func (i *Container) GetVersion() string {
	if i.ContainerSettings.Version == nil {
		return "latest"
	}
	return *i.ContainerSettings.Version
}

func (i *Container) GetImageNameWithTag() string {
	return *i.Service.Methods.Docker.Image + ":" + i.GetVersion()
}

func (i *Container) HasTag(tag string) bool {
	if i.ContainerSettings.Tags == nil {
		return false
	}
	for _, t := range i.ContainerSettings.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (i *Container) HasTagIn(tags []string) bool {
	if tags == nil {
		return true
	}
	for _, tag := range tags {
		if i.HasTag(tag) {
			return true
		}
	}
	return false
}
