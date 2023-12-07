package adapter

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	PathServices = "tests/services"
)

type AvailableTestSuite struct {
	suite.Suite

	adapter serviceFSAdapter
}

func TestAvailableTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableTestSuite))
}

func (suite *AvailableTestSuite) SetupSuite() {
	suite.adapter = *NewServiceFSAdapter(&ServiceFSAdapterParams{
		servicesPath: PathServices,
	}).(*serviceFSAdapter)

	err := suite.adapter.Reload()
	suite.Require().NoError(err)
	suite.NotZero(len(suite.adapter.services))
}

func (suite *AvailableTestSuite) TestGetAvailable() {
	suite.Len(suite.adapter.GetAll(), 1)
}
