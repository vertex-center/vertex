package types

import (
	"github.com/google/uuid"
)

const (
	EventNameInstancesChange = "change"

	EventNameInstanceStatusChange = "status_change"
	EventNameInstanceStdout       = "stdout"
	EventNameInstanceStderr       = "stderr"
)

type Listener interface {
	OnEvent(e interface{})
	GetUUID() uuid.UUID
}

type EventAdapterPort interface {
	AddListener(l Listener)
	RemoveListener(l Listener)
	Send(e interface{})
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

// Events

type EventInstanceLog struct {
	InstanceUUID uuid.UUID
	Kind         string
	Message      string
}

type EventInstanceStatusChange struct {
	InstanceUUID uuid.UUID
	Status       string
}

type EventInstancesChange struct{}
