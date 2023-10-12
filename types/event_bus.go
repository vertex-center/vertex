package types

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
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

	log.Debug("dispatching event", vlog.Any("count", len(*b.listeners)))

	// This code notify all listeners.
	// If some listeners are added while notifying, they will be
	// notified in the next loop, until all listeners are notified.

	notified := map[uuid.UUID]Listener{}
	tryCount := 0
	for len(notified) < len(*b.listeners) {
		if tryCount > 0 {
			log.Debug("some listeners were not notified; retrying...", vlog.Any("count", len(*b.listeners)-len(notified)))
		}
		if tryCount > 10 {
			log.Error(errors.New("after 10 retries to send events, there seems to be an issue with the event bus; the issue is probably caused by some listeners that create new listeners that themselves create new listeners, and so on"))
			break
		}

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

		for _, l := range toNotify {
			l.OnEvent(e)
			notified[l.GetUUID()] = l
		}
		tryCount++
	}
}
