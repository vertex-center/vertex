package event

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/event/types"
)

type TempListener struct {
	uuid    uuid.UUID
	onEvent func(e types.Event)
}

func NewTempListener(onEvent func(e types.Event)) TempListener {
	return TempListener{
		uuid:    uuid.New(),
		onEvent: onEvent,
	}
}

func (t TempListener) OnEvent(e types.Event) {
	t.onEvent(e)
}

func (t TempListener) GetUUID() uuid.UUID {
	return t.uuid
}
