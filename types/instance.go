package types

import (
	"encoding/json"
	"errors"
	"os/exec"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/nakabonne/tstorage"
	"github.com/vertex-center/vertex/logger"
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
	UseDocker   bool `json:"use_docker"`
	UseReleases bool `json:"use_releases"`

	// LaunchOnStartup indicates if the instance needs to start automatically when Vertex
	// starts. The default value is true.
	LaunchOnStartup *bool `json:"launch_on_startup,omitempty"`
}

type InstanceEvent struct {
	Name string
	Data string
}

type Instance struct {
	Service
	InstanceMetadata

	Status       string          `json:"status"`
	Logger       *InstanceLogger `json:"-"`
	EnvVariables EnvVariables    `json:"env"`

	UUID               uuid.UUID        `json:"uuid"`
	Cmd                *exec.Cmd        `json:"-"`
	UptimeStorage      tstorage.Storage `json:"-"`
	UptimeStopChannels []*chan bool     `json:"-"`

	Listeners map[uuid.UUID]chan InstanceEvent `json:"-"`
}

func NewInstance(id uuid.UUID, service Service, instancePath string) (Instance, error) {
	// TODO: Make UseDocker and UseReleases optional
	meta := InstanceMetadata{
		UseDocker:   false,
		UseReleases: false,
	}

	uptimeStorage, err := tstorage.NewStorage(
		tstorage.WithDataPath(path.Join(instancePath, ".vertex", "timestorage")),
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithWALBufferedSize(0),
	)
	if err != nil {
		return Instance{}, errors.New("failed to initialize time-storage")
	}

	return Instance{
		Service:          service,
		InstanceMetadata: meta,
		Status:           InstanceStatusOff,
		Logger:           NewInstanceLogger(instancePath),
		EnvVariables:     *NewEnvVariables(),
		UUID:             id,
		UptimeStorage:    uptimeStorage,
		Listeners:        map[uuid.UUID]chan InstanceEvent{},
	}, nil
}

type InstanceRepository interface {
	Get(uuid uuid.UUID) (*Instance, error)
	GetAll() map[uuid.UUID]*Instance
	GetPath(uuid uuid.UUID) string
	Delete(uuid uuid.UUID) error
	Exists(uuid uuid.UUID) bool
	Set(uuid uuid.UUID, instance Instance) error

	AddListener(channel chan InstanceEvent) uuid.UUID
	RemoveListener(uuid uuid.UUID)

	SaveMetadata(i *Instance) error
	LoadMetadata(i *Instance) error

	SaveEnv(i *Instance, variables map[string]string) error
	LoadEnv(i *Instance) error

	ReadService(instancePath string) (Service, error)

	Load(uuid uuid.UUID) (*Instance, error)

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

func (i *Instance) Register(channel chan InstanceEvent) uuid.UUID {
	id := uuid.New()
	i.Listeners[id] = channel

	logger.Log("registered to instance").
		AddKeyValue("channel", id).
		AddKeyValue("instance_uuid", i.UUID).
		Print()

	return id
}

func (i *Instance) Unregister(uuid uuid.UUID) {
	delete(i.Listeners, uuid)

	logger.Log("unregistered from instance").
		AddKeyValue("channel", uuid).
		AddKeyValue("instance_uuid", i.UUID).
		Print()
}

func (i *Instance) SetStatus(status string) {
	i.Status = status

	for _, listener := range i.Listeners {
		listener <- InstanceEvent{
			Name: InstanceEventStatusChange,
			Data: status,
		}
	}
}

func (i *Instance) NotifyListeners(event InstanceEvent) {
	for _, listener := range i.Listeners {
		listener <- event
	}
}

func (i *Instance) PushLogLine(line *LogLine) {
	i.Logger.Write(line)

	data, err := json.Marshal(line)
	if err != nil {
		logger.Error(err).Print()
	}

	var name string
	switch line.Kind {
	case InstanceEventStderr:
		name = InstanceEventStderr
	default:
		name = InstanceEventStdout
	}

	i.NotifyListeners(InstanceEvent{
		Name: name,
		Data: string(data),
	})
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

	i.NotifyListeners(InstanceEvent{
		Name: "uptime_status_change",
		Data: UptimeStatus(status),
	})

	return err
}
