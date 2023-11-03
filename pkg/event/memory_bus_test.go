package event

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/pkg/event/mock"
	"github.com/vertex-center/vertex/pkg/event/types"
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
	listener := mock.EventListener{}
	listener.OnEventFunc = func(e types.Event) {
		switch e.(type) {
		case mock.Event:
			called = true
		}
	}
	listener.GetUUIDFunc = func() uuid.UUID {
		return id
	}

	// Add a listener
	suite.bus.AddListener(&listener)
	assert.Equal(suite.T(), 1, len(*suite.bus.listeners))

	// Fire event
	suite.bus.DispatchEvent(mock.Event{})
	assert.Equal(suite.T(), 1, listener.OnEventCalls)
	assert.True(suite.T(), called)

	// Remove listener
	suite.bus.RemoveListener(&listener)
	assert.Equal(suite.T(), 0, len(*suite.bus.listeners))
}
