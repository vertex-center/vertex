package event

import (
	"errors"
	"sync"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vlog"
)

type MemoryBus struct {
	listeners      *map[uuid.UUID]Listener
	listenersMutex *sync.RWMutex
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		listeners:      &map[uuid.UUID]Listener{},
		listenersMutex: &sync.RWMutex{},
	}
}

func (b *MemoryBus) AddListener(l Listener) {
	b.listenersMutex.Lock()
	defer b.listenersMutex.Unlock()

	(*b.listeners)[l.GetUUID()] = l
}

func (b *MemoryBus) RemoveListener(l Listener) {
	b.listenersMutex.Lock()
	defer b.listenersMutex.Unlock()

	delete(*b.listeners, l.GetUUID())
}

func (b *MemoryBus) DispatchEvent(e Event) error {
	// This code notifies all listeners.
	// If some listeners are added while notifying, they will be
	// notified in the next loop, until all listeners are notified.

	notified := map[uuid.UUID]Listener{}
	tryCount := 0
	for {
		b.listenersMutex.RLock()
		var toNotify []Listener
		for _, l := range *b.listeners {
			if _, ok := notified[l.GetUUID()]; ok {
				// already notified
				continue
			}
			toNotify = append(toNotify, l)
		}
		b.listenersMutex.RUnlock()

		if len(toNotify) == 0 {
			// all listeners were notified
			break
		}

		if tryCount > 0 {
			log.Debug("some listeners were not notified; retrying...", vlog.Any("count", len(*b.listeners)-len(notified)))
		}
		if tryCount > 10 {
			log.Error(errors.New("after 10 retries to send events, there seems to be an issue with the event bus; the issue is probably caused by some listeners that create new listeners that themselves create new listeners, and so on"))
			break
		}

		for _, l := range toNotify {
			err := l.OnEvent(e)
			if err != nil {
				return err
			}
			notified[l.GetUUID()] = l
		}
		tryCount++
	}
	return nil
}
