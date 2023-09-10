package adapter

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventInMemoryAdapterTestSuite struct {
	suite.Suite

	adapter EventInMemoryAdapter
}

func TestEventInMemoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EventInMemoryAdapterTestSuite))
}

func (suite *EventInMemoryAdapterTestSuite) SetupSuite() {
	suite.adapter = *NewEventInMemoryAdapter().(*EventInMemoryAdapter)
}

func (suite *EventInMemoryAdapterTestSuite) TestEvents() {
	listener := MockListener{
		uuid: uuid.New(),
	}

	// Add a listener
	suite.adapter.AddListener(&listener)
	assert.Equal(suite.T(), 1, len(*suite.adapter.listeners))

	// Fire event
	listener.On("OnEvent").Return(nil)
	suite.adapter.Send(MockEvent{})
	listener.AssertCalled(suite.T(), "OnEvent")

	// Remove listener
	suite.adapter.RemoveListener(&listener)
	assert.Equal(suite.T(), 0, len(*suite.adapter.listeners))
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
