package adapter

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

type EventInMemoryAdapter struct {
	listeners *map[uuid.UUID]types.Listener
}

func NewEventInMemoryAdapter() types.EventAdapterPort {
	return &EventInMemoryAdapter{
		listeners: &map[uuid.UUID]types.Listener{},
	}
}

func (a *EventInMemoryAdapter) AddListener(l types.Listener) {
	(*a.listeners)[l.GetUUID()] = l
}

func (a *EventInMemoryAdapter) RemoveListener(l types.Listener) {
	delete(*a.listeners, l.GetUUID())
}

func (a *EventInMemoryAdapter) Send(e interface{}) {
	for _, l := range *a.listeners {
		l.OnEvent(e)
	}
}
