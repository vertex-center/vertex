package types

import (
	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
)

const (
	EventNameInstancesChange = "change"

	EventNameInstanceStatusChange = "status_change"
	EventNameInstanceStdout       = "stdout"
	EventNameInstanceStderr       = "stderr"
	EventNameInstanceDownload     = "download"
)

type Listener interface {
	OnEvent(e interface{})
	GetUUID() uuid.UUID
}

type TempListener struct {
	uuid    uuid.UUID
	onEvent func(e interface{})
}

func NewTempListener(onEvent func(e interface{})) TempListener {
	return TempListener{
		uuid:    uuid.New(),
		onEvent: onEvent,
	}
}

func (t TempListener) OnEvent(e interface{}) {
	t.onEvent(e)
}

func (t TempListener) GetUUID() uuid.UUID {
	return t.uuid
}

type (
	EventServerStart         struct{}
	EventServerStop          struct{}
	EventDependenciesUpdated struct{}
)

// Events instance

type EventInstanceLoaded struct {
	Instance *instancestypes.Instance
}

type EventInstanceLog struct {
	InstanceUUID uuid.UUID
	Kind         string
	Message      instancestypes.LogLineMessage
}

type EventInstanceStatusChange struct {
	InstanceUUID uuid.UUID
	ServiceID    string
	Instance     instancestypes.Instance
	Name         string
	Status       string
}

type EventInstanceCreated struct{}
type EventInstanceDeleted struct {
	InstanceUUID uuid.UUID
	ServiceID    string
}

// Events instances

type EventInstancesChange struct{}

type EventInstancesLoaded struct {
	Count int
}
