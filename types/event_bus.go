package types

import (
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
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
	if _, ok := e.(EventServerHardReset); ok {
		if !config.Current.Debug() {
			log.Warn("hard reset event received but skipped; this can be a malicious application, or you may have forgotten to switch to the development mode.")
			return
		}
		log.Warn("hard reset event dispatched.")
	}

	b.listenersMutex.RLock()
	defer b.listenersMutex.RUnlock()

	for _, l := range *b.listeners {
		l.OnEvent(e)
	}
}
