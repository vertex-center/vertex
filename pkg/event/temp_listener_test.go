package event

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TempListenerTestSuite struct {
	suite.Suite
}

func TestTempListenerTestSuite(t *testing.T) {
	suite.Run(t, new(TempListenerTestSuite))
}

func (suite *TempListenerTestSuite) TestOnEvent() {
	called := false
	listener := NewTempListener(func(e Event) {
		switch e.(type) {
		case MockEvent:
			called = true
		}
	})
	listener.OnEvent(struct{}{})
	suite.False(called)
	listener.OnEvent(MockEvent{})
	suite.True(called)
}

func (suite *TempListenerTestSuite) TestGetUUID() {
	listener := NewTempListener(func(e Event) {})
	suite.NotEqual(uuid.Nil, listener.GetUUID())
}
