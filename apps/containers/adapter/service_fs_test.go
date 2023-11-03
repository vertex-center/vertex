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

	adapter ServiceFSAdapter
}

func TestAvailableTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableTestSuite))
}

func (suite *AvailableTestSuite) SetupSuite() {
	suite.adapter = *NewServiceFSAdapter(&ServiceFSAdapterParams{
		servicesPath: PathServices,
	}).(*ServiceFSAdapter)

	err := suite.adapter.Reload()
	suite.NoError(err)
	suite.NotZero(len(suite.adapter.services))
}

func (suite *AvailableTestSuite) TestGetAvailable() {
	suite.Equal(1, len(suite.adapter.GetAll()))
}
