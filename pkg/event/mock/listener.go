package mock

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/event/types"
)

type EventListener struct {
	OnEventFunc  func(e types.Event)
	OnEventCalls int
	GetUUIDFunc  func() uuid.UUID
	GetUUIDCalls int
}

func (m *EventListener) OnEvent(e types.Event) {
	m.OnEventCalls++
	m.OnEventFunc(e)
}

func (m *EventListener) GetUUID() uuid.UUID {
	m.GetUUIDCalls++
	return m.GetUUIDFunc()
}
