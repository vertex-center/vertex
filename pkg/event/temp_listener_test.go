package event

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/common/uuid"
)

type TempListenerTestSuite struct {
	suite.Suite
}

func TestTempListenerTestSuite(t *testing.T) {
	suite.Run(t, new(TempListenerTestSuite))
}

func (suite *TempListenerTestSuite) TestOnEvent() {
	called := false
	listener := NewTempListener(func(e Event) error {
		switch e.(type) {
		case MockEvent:
			called = true
		}
		return nil
	})

	err := listener.OnEvent(struct{}{})
	suite.Require().NoError(err)
	suite.False(called)

	err = listener.OnEvent(MockEvent{})
	suite.Require().NoError(err)
	suite.True(called)
}

func (suite *TempListenerTestSuite) TestGetUUID() {
	listener := NewTempListener(func(e Event) error {
		return nil
	})
	suite.NotEqual(uuid.Nil, listener.GetUUID())
}
