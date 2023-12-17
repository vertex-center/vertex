package event

import (
	"github.com/vertex-center/vertex/common/uuid"
)

type TempListener struct {
	uuid    uuid.UUID
	onEvent func(e Event) error
}

func NewTempListener(onEvent func(e Event) error) TempListener {
	return TempListener{
		uuid:    uuid.New(),
		onEvent: onEvent,
	}
}

func (t TempListener) OnEvent(e Event) error {
	return t.onEvent(e)
}

func (t TempListener) GetUUID() uuid.UUID {
	return t.uuid
}
