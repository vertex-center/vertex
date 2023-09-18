package adapter

import (
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

type EventInMemoryAdapter struct {
	listeners      *map[uuid.UUID]types.Listener
	listenersMutex *sync.RWMutex
}

func NewEventInMemoryAdapter() types.EventAdapterPort {
	return &EventInMemoryAdapter{
		listeners:      &map[uuid.UUID]types.Listener{},
		listenersMutex: &sync.RWMutex{},
	}
}

func (a *EventInMemoryAdapter) AddListener(l types.Listener) {
	a.listenersMutex.Lock()
	defer a.listenersMutex.Unlock()

	(*a.listeners)[l.GetUUID()] = l
}

func (a *EventInMemoryAdapter) RemoveListener(l types.Listener) {
	a.listenersMutex.Lock()
	defer a.listenersMutex.Unlock()

	delete(*a.listeners, l.GetUUID())
}

func (a *EventInMemoryAdapter) Send(e interface{}) {
	a.listenersMutex.RLock()
	defer a.listenersMutex.RUnlock()

	for _, l := range *a.listeners {
		l.OnEvent(e)
	}
}
