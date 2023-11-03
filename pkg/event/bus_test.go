package event

type MockBus struct {
	AddListenerFunc     func(l Listener)
	AddListenerCalls    int
	RemoveListenerFunc  func(l Listener)
	RemoveListenerCalls int
	DispatchEventFunc   func(e Event)
	DispatchEventCalls  int
}

func (m *MockBus) AddListener(l Listener) {
	m.AddListenerCalls++
	m.AddListenerFunc(l)
}

func (m *MockBus) RemoveListener(l Listener) {
	m.RemoveListenerCalls++
	m.RemoveListenerFunc(l)
}

func (m *MockBus) DispatchEvent(e Event) {
	m.DispatchEventCalls++
	m.DispatchEventFunc(e)
}
