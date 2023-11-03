package event

import (
	"github.com/google/uuid"
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
