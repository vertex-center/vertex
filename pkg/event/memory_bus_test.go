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

	called := false
	listener := &MockListener{}
	listener.OnEventFunc = func(e Event) {
		switch e.(type) {
		case MockEvent:
			called = true
		}
	}
	listener.GetUUIDFunc = func() uuid.UUID {
		return id
	}

	// Add a listener
	suite.bus.AddListener(listener)
	suite.Len(*suite.bus.listeners, 1)

	// Fire event
	suite.bus.DispatchEvent(MockEvent{})
	suite.Equal(1, listener.OnEventCalls)
	suite.True(called)

	// Remove listener
	suite.bus.RemoveListener(listener)
	suite.Empty(*suite.bus.listeners)
}
