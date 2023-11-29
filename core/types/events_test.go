package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EventsTestSuite struct {
	suite.Suite
}

func TestEventsTestSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}

type (
	MockObjA struct{}
	MockObjB struct{}
)

func (suite *EventsTestSuite) TestEventDbCopy() {
	event := NewEventDbCopy()
	suite.NotNil(event.tables)
	suite.Empty(event.All())

	event.AddTable(MockObjA{})
	event.AddTable(MockObjB{})

	suite.NotEmpty(event.All())

	expected := []interface{}{
		MockObjA{},
		MockObjB{},
	}
	suite.Equal(expected, event.All())
}
