package servicesmanager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	PathLive     = "live_test"
	PathServices = "live_test/services"
)

type AvailableTestSuite struct {
	suite.Suite
}

func TestAvailableTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableTestSuite))
}

func (suite *AvailableTestSuite) SetupSuite() {
	testReload(suite.T())
}

func (suite *AvailableTestSuite) TearDownSuite() {
	err := os.RemoveAll(PathLive)
	assert.NoError(suite.T(), err)
}

func testReload(t *testing.T) {
	err := os.MkdirAll(PathServices, os.ModePerm)
	assert.NoError(t, err)

	err = reload(PathServices)
	assert.NoError(t, err)

	assert.NotZero(t, len(available))
}

func (suite *AvailableTestSuite) TestListAvailable() {
	assert.NotZero(suite.T(), len(ListAvailable()))
}
