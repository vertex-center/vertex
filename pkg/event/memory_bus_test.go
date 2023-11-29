package event

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MemoryBusTestSuite struct {
	suite.Suite

	bus MemoryBus
}

func TestMemoryBusTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryBusTestSuite))
}

func (suite *MemoryBusTestSuite) SetupSuite() {
	suite.bus = *NewMemoryBus()
}

func (suite *MemoryBusTestSuite) TestEvents() {
	id := uuid.New()

	listener := &MockListener{}
	listener.On("OnEvent", MockEvent{}).Return()
	listener.On("GetUUID").Return(id)

	// Add a listener
	suite.bus.AddListener(listener)
	suite.Len(*suite.bus.listeners, 1)

	// Fire event
	err := suite.bus.DispatchEvent(MockEvent{})
	suite.Require().NoError(err)
	listener.AssertExpectations(suite.T())

	// Remove listener
	suite.bus.RemoveListener(listener)
	suite.Empty(*suite.bus.listeners)
}
