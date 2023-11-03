package event

import (
	"github.com/google/uuid"
)

type TempListener struct {
	uuid    uuid.UUID
	onEvent func(e Event)
}

func NewTempListener(onEvent func(e Event)) TempListener {
	return TempListener{
		uuid:    uuid.New(),
		onEvent: onEvent,
	}
}

func (t TempListener) OnEvent(e Event) {
	t.onEvent(e)
}

func (t TempListener) GetUUID() uuid.UUID {
	return t.uuid
}
