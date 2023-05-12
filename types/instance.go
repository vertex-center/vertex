package types

import (
	"errors"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/nakabonne/tstorage"
)

const (
	InstanceStatusOff      = "off"
	InstanceStatusBuilding = "building"
	InstanceStatusStarting = "starting"
	InstanceStatusRunning  = "running"
	InstanceStatusError    = "error"

	InstanceEventStatusChange = "status_change"
	InstanceEventStdout       = "stdout"
	InstanceEventStderr       = "stderr"
)

type InstanceMetadata struct {
	// UseDocker indicates if the instance should be launched with Docker.
	// The default value is false.
	UseDocker *bool `json:"use_docker,omitempty"`

	// UseReleases indicates if the instance should use precompiled releases when possible.
	// The default value is false.
	UseReleases *bool `json:"use_releases,omitempty"`

	// LaunchOnStartup indicates if the instance needs to start automatically when Vertex starts.
	// The default value is true.
	LaunchOnStartup *bool `json:"launch_on_startup,omitempty"`
}

type InstanceEvent struct {
	Name string
	Data string
}

type EnvVariables map[string]string

type Instance struct {
	Service
	InstanceMetadata

	Status       string       `json:"status"`
	EnvVariables EnvVariables `json:"env"`

	UUID               uuid.UUID        `json:"uuid"`
	UptimeStorage      tstorage.Storage `json:"-"`
	UptimeStopChannels []*chan bool     `json:"-"`
}

func NewInstance(id uuid.UUID, service Service, instancePath string) (Instance, error) {
	uptimeStorage, err := tstorage.NewStorage(
		tstorage.WithDataPath(path.Join(instancePath, ".vertex", "timestorage")),
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithWALBufferedSize(0),
	)
	if err != nil {
		return Instance{}, errors.New("failed to initialize time-storage")
	}

	return Instance{
		Service:       service,
		Status:        InstanceStatusOff,
		EnvVariables:  map[string]string{},
		UUID:          id,
		UptimeStorage: uptimeStorage,
	}, nil
}

type InstanceRepository interface {
	Get(uuid uuid.UUID) (*Instance, error)
	GetAll() map[uuid.UUID]*Instance
	GetPath(uuid uuid.UUID) string
	Delete(uuid uuid.UUID) error
	Exists(uuid uuid.UUID) bool
	Set(uuid uuid.UUID, instance Instance) error

	SaveMetadata(i *Instance) error
	LoadMetadata(i *Instance) error

	SaveEnv(i *Instance, variables map[string]string) error
	LoadEnv(i *Instance) error

	ReadService(instancePath string) (Service, error)

	Reload(func(uuid uuid.UUID))

	Close()
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

func (i *Instance) PushStatus(name string, status float64) error {
	err := i.UptimeStorage.InsertRows([]tstorage.Row{
		{
			Metric: "status_change",
			Labels: []tstorage.Label{
				{
					Name:  "name",
					Value: name,
				},
			},
			DataPoint: tstorage.DataPoint{
				Timestamp: time.Now().Unix(),
				Value:     status,
			},
		},
	})

	//i.NotifyListeners(InstanceEvent{
	//	Name: "uptime_status_change",
	//	Data: UptimeStatus(status),
	//})

	return err
}

func (i *Instance) IsDockerized() bool {
	return i.InstanceMetadata.UseDocker != nil && *i.InstanceMetadata.UseDocker
}
