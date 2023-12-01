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

func (suite *EventsTestSuite) TestEventDbCopy() {
	event := NewEventDbCopy()
	suite.NotNil(event.tables)
	suite.Empty(event.All())

	event.AddTable("table_a")
	event.AddTable("table_b")

	suite.NotEmpty(event.All())
	suite.Equal([]string{
		"table_a",
		"table_b",
	}, event.All())
}
