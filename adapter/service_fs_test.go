package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(suite.T(), err)

	assert.NotZero(suite.T(), len(suite.adapter.services))
}

func (suite *AvailableTestSuite) TestGetAvailable() {
	assert.Equal(suite.T(), 1, len(suite.adapter.GetAll()))
}
