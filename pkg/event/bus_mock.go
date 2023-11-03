package event

type MockBus struct {
	AddListenerFunc     func(l EventListener)
	AddListenerCalls    int
	RemoveListenerFunc  func(l EventListener)
	RemoveListenerCalls int
	DispatchEventFunc   func(e Event)
	DispatchEventCalls  int
}

func (m *MockBus) AddListener(l EventListener) {
	m.AddListenerCalls++
	m.AddListenerFunc(l)
}

func (m *MockBus) RemoveListener(l EventListener) {
	m.RemoveListenerCalls++
	m.RemoveListenerFunc(l)
}

func (m *MockBus) DispatchEvent(e Event) {
	m.DispatchEventCalls++
	m.DispatchEventFunc(e)
}
