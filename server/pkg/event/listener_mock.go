package event

import (
	"github.com/stretchr/testify/mock"
	"github.com/vertex-center/uuid"
)

type MockListener struct{ mock.Mock }

var _ Listener = (*MockListener)(nil)

func (m *MockListener) OnEvent(e Event) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockListener) GetUUID() uuid.UUID {
	args := m.Called()
	return args.Get(0).(uuid.UUID)
}
