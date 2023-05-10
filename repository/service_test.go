package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	PathServices = "services_test_live"
)

type AvailableTestSuite struct {
	suite.Suite

	repo ServiceRepository
}

func TestAvailableTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableTestSuite))
}

func (suite *AvailableTestSuite) SetupSuite() {
	suite.repo = NewServiceRepository(&ServiceRepositoryParams{
		servicesPath: PathServices,
	})

	err := os.MkdirAll(PathServices, os.ModePerm)
	assert.NoError(suite.T(), err)

	err = suite.repo.reload()
	assert.NoError(suite.T(), err)

	assert.NotZero(suite.T(), len(suite.repo.services))
}

func (suite *AvailableTestSuite) TearDownSuite() {
	err := os.RemoveAll(PathServices)
	assert.NoError(suite.T(), err)
}

func (suite *AvailableTestSuite) TestGetAvailable() {
	assert.NotZero(suite.T(), suite.repo.GetAll())
}
