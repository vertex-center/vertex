package mock

import (
	"github.com/vertex-center/vertex/pkg/event/types"
)

type EventBus struct {
	AddListenerFunc     func(l types.EventListener)
	AddListenerCalls    int
	RemoveListenerFunc  func(l types.EventListener)
	RemoveListenerCalls int
	DispatchEventFunc   func(e types.Event)
	DispatchEventCalls  int
}

func (m EventBus) AddListener(l types.EventListener) {
	m.AddListenerCalls++
	m.AddListenerFunc(l)
}

func (m EventBus) RemoveListener(l types.EventListener) {
	m.RemoveListenerCalls++
	m.RemoveListenerFunc(l)
}

func (m EventBus) DispatchEvent(e types.Event) {
	m.DispatchEventCalls++
	m.DispatchEventFunc(e)
}
