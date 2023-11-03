package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/pkg/event"
)

type VertexContextTestSuite struct {
	suite.Suite

	context *VertexContext
}

func TestVertexContextTestSuite(t *testing.T) {
	suite.Run(t, new(VertexContextTestSuite))
}

func (suite *VertexContextTestSuite) SetupSuite() {
	suite.context = NewVertexContext()
	suite.NotNil(suite.context.bus)
}

func (suite *VertexContextTestSuite) TestDispatchEvent() {
	bus := &event.MockBus{}
	bus.DispatchEventFunc = func(e event.Event) {}
	suite.context.bus = bus

	suite.context.DispatchEvent(event.MockEvent{})
	suite.Equal(1, bus.DispatchEventCalls)
}

func (suite *VertexContextTestSuite) TestDispatchHardReset() {
	bus := &event.MockBus{}
	suite.context.bus = bus

	suite.context.DispatchEvent(EventServerHardReset{})
	suite.Equal(0, bus.DispatchEventCalls)
}

func (suite *VertexContextTestSuite) TestAddListener() {
	bus := &event.MockBus{}
	bus.AddListenerFunc = func(l event.EventListener) {}
	suite.context.bus = bus

	suite.context.AddListener(&event.MockListener{})
	suite.Equal(1, bus.AddListenerCalls)
}

func (suite *VertexContextTestSuite) TestRemoveListener() {
	bus := &event.MockBus{}
	bus.RemoveListenerFunc = func(l event.EventListener) {}
	suite.context.bus = bus

	listener := &event.MockListener{}
	suite.context.RemoveListener(listener)
	suite.Equal(1, bus.RemoveListenerCalls)
}
