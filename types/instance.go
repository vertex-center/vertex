package types

import (
	"errors"

	"github.com/google/uuid"
)

const (
	InstanceStatusOff      = "off"
	InstanceStatusBuilding = "building"
	InstanceStatusStarting = "starting"
	InstanceStatusRunning  = "running"
	InstanceStatusStopping = "stopping"
	InstanceStatusError    = "error"
)

const (
	InstanceInstallMethodScript  = "script"
	InstanceInstallMethodRelease = "release"
	InstanceInstallMethodDocker  = "docker"
)

var (
	ErrInstanceNotFound     = errors.New("instance not found")
	ErrInstanceStillRunning = errors.New("instance still running")
)

type Instance struct {
	InstanceSettings

	Service Service              `json:"service"`
	UUID    uuid.UUID            `json:"uuid"`
	Status  string               `json:"status"`
	Env     InstanceEnvVariables `json:"environment,omitempty"`

	Update        *InstanceUpdate `json:"update,omitempty"`
	ServiceUpdate ServiceUpdate   `json:"service_update,omitempty"`

	CacheVersions []string `json:"cache_versions,omitempty"`
}

type InstanceQuery struct {
	Features []string `json:"features,omitempty"`
}

type InstanceUpdate struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

type DownloadProgress struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Current int64  `json:"current,omitempty"`
	Total   int64  `json:"total,omitempty"`
}

func NewInstance(id uuid.UUID, service Service) Instance {
	return Instance{
		Service: service,
		UUID:    id,
		Status:  InstanceStatusOff,
		Env:     map[string]string{},
	}
}

type InstanceAdapterPort interface {
	GetPath(uuid uuid.UUID) string
	Delete(uuid uuid.UUID) error

	GetAll() ([]uuid.UUID, error)
}

func (i *Instance) DockerImageVertexName() string {
	return "vertex_image_" + i.UUID.String()
}

func (i *Instance) DockerContainerName() string {
	return "VERTEX_CONTAINER_" + i.UUID.String()
}

func (i *Instance) IsRunning() bool {
	return i.Status != InstanceStatusOff && i.Status != InstanceStatusError
}

func (i *Instance) IsBusy() bool {
	return i.Status == InstanceStatusBuilding || i.Status == InstanceStatusStarting || i.Status == InstanceStatusStopping
}

func (i *Instance) IsDockerized() bool {
	return i.InstanceSettings.InstallMethod != nil && *i.InstanceSettings.InstallMethod == InstanceInstallMethodDocker
}

func (i *Instance) LaunchOnStartup() bool {
	launchOnStartup := i.InstanceSettings.LaunchOnStartup
	if launchOnStartup != nil && !*launchOnStartup {
		return false
	}
	return true
}

func (i *Instance) ResetDefaultEnv() {
	i.Env = InstanceEnvVariables{}
	for _, env := range i.Service.Env {
		i.Env[env.Name] = env.Default
	}
}

func (i *Instance) HasFeature(featureType string) bool {
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

func (i *Instance) HasOneOfFeatures(featureTypes []string) bool {
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

func (i *Instance) GetVersion() string {
	if i.InstanceSettings.Version == nil {
		return "latest"
	}
	return *i.InstanceSettings.Version
}

func (i *Instance) GetImageNameWithTag() string {
	return *i.Service.Methods.Docker.Image + ":" + i.GetVersion()
}
