package event

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockListener struct{ mock.Mock }

func (m *MockListener) OnEvent(e Event) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockListener) GetUUID() uuid.UUID {
	args := m.Called()
	return args.Get(0).(uuid.UUID)
}
