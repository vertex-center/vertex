package event

import "github.com/stretchr/testify/mock"

type MockBus struct {
	mock.Mock
}

func (m *MockBus) AddListener(l Listener) {
	m.Called(l)
}

func (m *MockBus) RemoveListener(l Listener) {
	m.Called(l)
}

func (m *MockBus) DispatchEvent(e Event) error {
	args := m.Called(e)
	return args.Error(0)
}
