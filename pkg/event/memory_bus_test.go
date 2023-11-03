package event

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MemoryBusTestSuite struct {
	suite.Suite

	bus MemoryBus
}

func TestEventInMemoryAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryBusTestSuite))
}

func (suite *MemoryBusTestSuite) SetupSuite() {
	suite.bus = *NewMemoryBus()
}

func (suite *MemoryBusTestSuite) TestEvents() {
	listener := MockListener{
		uuid: uuid.New(),
	}

	// Add a listener
	suite.bus.AddListener(&listener)
	assert.Equal(suite.T(), 1, len(*suite.bus.listeners))

	// Fire event
	listener.On("OnEvent").Return(nil)
	suite.bus.DispatchEvent(MockEvent{})
	listener.AssertCalled(suite.T(), "OnEvent")

	// Remove listener
	suite.bus.RemoveListener(&listener)
	assert.Equal(suite.T(), 0, len(*suite.bus.listeners))
}

type MockEvent struct{}

type MockListener struct {
	mock.Mock

	uuid uuid.UUID
}

func (m *MockListener) OnEvent(e interface{}) {
	switch e.(type) {
	case MockEvent:
		m.Called()
	}
}

func (m *MockListener) GetUUID() uuid.UUID {
	return m.uuid
}
