package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const PathLive = "live_test"
const PathRepo = "live_test/repo"

type RepositoryTestSuite struct {
	suite.Suite
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	err := os.RemoveAll(PathLive)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestDownloadLatestRepository() {
	err := os.MkdirAll(PathRepo, os.ModePerm)
	assert.NoError(suite.T(), err)

	url := "https://github.com/vertex-center/vertex-services"

	// reload to test Clone()
	err = DownloadLatestRepository(PathRepo, url)
	assert.NoError(suite.T(), err)
	assert.DirExists(suite.T(), PathRepo)

	// reload to test Pull()
	err = DownloadLatestRepository(PathRepo, url)
	assert.NoError(suite.T(), err)
	assert.DirExists(suite.T(), PathRepo)
}
