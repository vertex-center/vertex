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

	adapter templateFSAdapter
}

func TestAvailableTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableTestSuite))
}

func (suite *AvailableTestSuite) SetupSuite() {
	suite.adapter = *NewTemplateFSAdapter(&TemplateFSAdapterParams{
		templatesPath: PathServices,
	}).(*templateFSAdapter)

	err := suite.adapter.Reload()
	suite.Require().NoError(err)
	suite.NotZero(len(suite.adapter.templates))
}

func (suite *AvailableTestSuite) TestGetAvailable() {
	suite.Len(suite.adapter.GetAll(), 1)
}
