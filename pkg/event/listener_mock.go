package event

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockListener struct {
	mock.Mock
}

func (m *MockListener) OnEvent(e Event) {
	m.Called(e)
}

func (m *MockListener) GetUUID() uuid.UUID {
	args := m.Called()
	return args.Get(0).(uuid.UUID)
}
