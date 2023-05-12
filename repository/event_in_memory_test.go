package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventInMemoryRepositoryTestSuite struct {
	suite.Suite

	repo EventInMemoryRepository
}

func TestEventInMemoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EventInMemoryRepositoryTestSuite))
}

func (suite *EventInMemoryRepositoryTestSuite) SetupSuite() {
	suite.repo = NewEventInMemoryRepository()
}

func (suite *EventInMemoryRepositoryTestSuite) TestEvents() {
	listener := MockListener{
		uuid: uuid.New(),
	}

	// Add a listener
	suite.repo.AddListener(&listener)
	assert.Equal(suite.T(), 1, len(*suite.repo.listeners))

	// Fire event
	listener.On("OnEvent").Return(nil)
	suite.repo.Send(MockEvent{})
	listener.AssertCalled(suite.T(), "OnEvent")

	// Remove listener
	suite.repo.RemoveListener(&listener)
	assert.Equal(suite.T(), 0, len(*suite.repo.listeners))
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
