package event

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/pkg/event/mock"
	"github.com/vertex-center/vertex/pkg/event/types"
)

type TempListenerTestSuite struct {
	suite.Suite
}

func TestTempListenerTestSuite(t *testing.T) {
	suite.Run(t, new(TempListenerTestSuite))
}

func (suite *TempListenerTestSuite) TestOnEvent() {
	called := false
	listener := NewTempListener(func(e types.Event) {
		switch e.(type) {
		case mock.Event:
			called = true
		}
	})
	listener.OnEvent(struct{}{})
	suite.False(called)
	listener.OnEvent(mock.Event{})
	suite.True(called)
}

func (suite *TempListenerTestSuite) TestGetUUID() {
	listener := NewTempListener(func(e types.Event) {})
	suite.NotEqual(uuid.Nil, listener.GetUUID())
}
