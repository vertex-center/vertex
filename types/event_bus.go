package types

import (
	"sync"

	"github.com/google/uuid"
)

type EventBus struct {
	listeners      *map[uuid.UUID]Listener
	listenersMutex *sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners:      &map[uuid.UUID]Listener{},
		listenersMutex: &sync.RWMutex{},
	}
}

func (b *EventBus) AddListener(l Listener) {
	b.listenersMutex.Lock()
	defer b.listenersMutex.Unlock()

	(*b.listeners)[l.GetUUID()] = l
}

func (b *EventBus) RemoveListener(l Listener) {
	b.listenersMutex.Lock()
	defer b.listenersMutex.Unlock()

	delete(*b.listeners, l.GetUUID())
}

func (b *EventBus) Send(e interface{}) {
	b.listenersMutex.RLock()
	defer b.listenersMutex.RUnlock()

	for _, l := range *b.listeners {
		l.OnEvent(e)
	}
}
