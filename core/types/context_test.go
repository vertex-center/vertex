package types

import (
	"testing"

	"github.com/stretchr/testify/mock"
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
	bus.On("DispatchEvent", event.MockEvent{}).Return()
	suite.context.bus = bus

	suite.context.DispatchEvent(event.MockEvent{})
	bus.AssertExpectations(suite.T())
}

func (suite *VertexContextTestSuite) TestDispatchHardReset() {
	bus := &event.MockBus{}
	suite.context.bus = bus

	suite.context.DispatchEvent(EventServerHardReset{})
	bus.AssertNotCalled(suite.T(), "DispatchEvent", mock.Anything)
}

func (suite *VertexContextTestSuite) TestAddListener() {
	bus := &event.MockBus{}
	bus.On("AddListener", &event.MockListener{}).Return()
	suite.context.bus = bus

	suite.context.AddListener(&event.MockListener{})
	bus.AssertExpectations(suite.T())
}

func (suite *VertexContextTestSuite) TestRemoveListener() {
	bus := &event.MockBus{}
	bus.On("RemoveListener", &event.MockListener{}).Return()
	suite.context.bus = bus

	listener := &event.MockListener{}
	suite.context.RemoveListener(listener)
	bus.AssertExpectations(suite.T())
}