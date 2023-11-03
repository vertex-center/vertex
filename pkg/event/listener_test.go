package event

import "github.com/google/uuid"

type MockListener struct {
	OnEventFunc  func(e Event)
	OnEventCalls int
	GetUUIDFunc  func() uuid.UUID
	GetUUIDCalls int
}

func (m *MockListener) OnEvent(e Event) {
	m.OnEventCalls++
	m.OnEventFunc(e)
}

func (m *MockListener) GetUUID() uuid.UUID {
	m.GetUUIDCalls++
	return m.GetUUIDFunc()
}
