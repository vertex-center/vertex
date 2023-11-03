package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/pkg/event/mock"
	"github.com/vertex-center/vertex/pkg/event/types"
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
	bus := &mock.EventBus{}
	bus.DispatchEventFunc = func(e types.Event) {}
	suite.context.bus = bus

	suite.context.DispatchEvent(mock.Event{})
	suite.Equal(1, bus.DispatchEventCalls)
}

func (suite *VertexContextTestSuite) TestDispatchHardReset() {
	bus := &mock.EventBus{}
	suite.context.bus = bus

	suite.context.DispatchEvent(EventServerHardReset{})
	suite.Equal(0, bus.DispatchEventCalls)
}

func (suite *VertexContextTestSuite) TestAddListener() {
	bus := &mock.EventBus{}
	bus.AddListenerFunc = func(l types.EventListener) {}
	suite.context.bus = bus

	suite.context.AddListener(&mock.EventListener{})
	suite.Equal(1, bus.AddListenerCalls)
}

func (suite *VertexContextTestSuite) TestRemoveListener() {
	bus := &mock.EventBus{}
	bus.RemoveListenerFunc = func(l types.EventListener) {}
	suite.context.bus = bus

	listener := &mock.EventListener{}
	suite.context.RemoveListener(listener)
	suite.Equal(1, bus.RemoveListenerCalls)
}
