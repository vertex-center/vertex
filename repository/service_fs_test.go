package repository

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

	repo ServiceFSRepository
}

func TestAvailableTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableTestSuite))
}

func (suite *AvailableTestSuite) SetupSuite() {
	suite.repo = NewServiceFSRepository(&ServiceRepositoryParams{
		servicesPath: PathServices,
	})

	err := suite.repo.reload()
	assert.NoError(suite.T(), err)

	assert.NotZero(suite.T(), len(suite.repo.services))
}

func (suite *AvailableTestSuite) TestGetAvailable() {
	assert.Equal(suite.T(), 1, len(suite.repo.GetAll()))
}
