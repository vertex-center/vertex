package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/core/types"
)

type DebugServiceTestSuite struct {
	suite.Suite

	service *DebugService
}

func (suite *DebugServiceTestSuite) SetupSuite() {
	ctx := types.NewVertexContext(&types.DB{})
	suite.service = NewDebugService(ctx).(*DebugService)
}

func TestDebugServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DebugServiceTestSuite))
}

func (suite *DebugServiceTestSuite) TestHardReset() {
	// TODO: test if the event is dispatched
	// this will require some rework of the event system
	// to allow for mocking.
	suite.service.HardReset()
}
