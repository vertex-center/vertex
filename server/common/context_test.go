package common

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/server/pkg/event"
)

type VertexContextTestSuite struct {
	suite.Suite

	context *VertexContext
}

func TestVertexContextTestSuite(t *testing.T) {
	suite.Run(t, new(VertexContextTestSuite))
}

func (suite *VertexContextTestSuite) SetupSuite() {
	suite.context = NewVertexContext(About{}, false)
	suite.NotNil(suite.context.bus)
}

func (suite *VertexContextTestSuite) TestDispatchEvent() {
	bus := &event.MockBus{}
	bus.On("DispatchEvent", event.MockEvent{}).Return(nil)
	suite.context.bus = bus

	err := suite.context.DispatchEventWithErr(event.MockEvent{})
	suite.Require().NoError(err)
	bus.AssertExpectations(suite.T())
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
