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

type InstanceSettings struct {
	// Method indicates how the instance is installed.
	// It can be by script, release or docker.
	InstallMethod *string `json:"install_method,omitempty"`

	// LaunchOnStartup indicates if the instance needs to start automatically when Vertex starts.
	// The default value is true.
	LaunchOnStartup *bool `json:"launch_on_startup,omitempty"`

	// DisplayName is a custom name for the instance.
	DisplayName *string `json:"display_name,omitempty"`
}

type EnvVariables map[string]string

type Instance struct {
	InstanceSettings

	Service Service      `json:"service"`
	UUID    uuid.UUID    `json:"uuid"`
	Status  string       `json:"status"`
	Env     EnvVariables `json:"environment,omitempty"`

	Update *InstanceUpdate `json:"update,omitempty"`
}

type InstanceQuery struct {
	Features []string `json:"features,omitempty"`
}

type InstanceUpdate struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
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
	Get(uuid uuid.UUID) (*Instance, error)
	GetAll() map[uuid.UUID]*Instance
	Search(query InstanceQuery) map[uuid.UUID]*Instance
	GetPath(uuid uuid.UUID) string
	Delete(uuid uuid.UUID) error
	Exists(uuid uuid.UUID) bool
	Set(uuid uuid.UUID, instance Instance) error

	SaveSettings(i *Instance) error
	LoadSettings(i *Instance) error

	SaveEnv(i *Instance, variables map[string]string) error
	LoadEnv(i *Instance) error

	ReadService(instancePath string) (Service, error)

	Reload(func(uuid uuid.UUID))
}

func (i *Instance) DockerImageName() string {
	return "vertex_image_" + i.UUID.String()
}

func (i *Instance) DockerContainerName() string {
	return "VERTEX_CONTAINER_" + i.UUID.String()
}

func (i *Instance) IsRunning() bool {
	return i.Status != InstanceStatusOff && i.Status != InstanceStatusError
}

func (i *Instance) IsDockerized() bool {
	return i.InstanceSettings.InstallMethod != nil && *i.InstanceSettings.InstallMethod == InstanceInstallMethodDocker
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
