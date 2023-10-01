package migration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationTestSuite struct {
	suite.Suite

	tool *MigrationTool
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}

func (suite *MigrationTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "live_temp-*")
	if err != nil {
		return
	}

	suite.tool = NewMigrationTool(dir)
}

func (suite *MigrationTestSuite) TearDownTest() {
	err := os.RemoveAll(suite.tool.livePath)
	suite.NoError(err)
}

func (suite *MigrationTestSuite) TestMigrate() {
	err := suite.tool.Migrate()
	suite.NoError(err)

	_, err = os.Stat(suite.tool.metadataPath)
	suite.NoError(err)

	v, err := suite.tool.readLiveVersion()
	suite.NoError(err)
	suite.Equal(len(suite.tool.migrations)-1, v.Version)
}
