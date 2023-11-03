package migration

import (
	"os"
	"path"
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
	suite.Require().NoError(err)

	err = os.MkdirAll(path.Join(dir, "instances"), 0755)
	suite.Require().NoError(err)

	suite.tool = NewMigrationTool(dir)
}

func (suite *MigrationTestSuite) TearDownTest() {
	err := os.RemoveAll(suite.tool.livePath)
	suite.Require().NoError(err)
}

func (suite *MigrationTestSuite) TestMigrate() {
	_, err := suite.tool.Migrate()
	suite.Require().NoError(err)

	_, err = os.Stat(suite.tool.metadataPath)
	suite.Require().NoError(err)

	v, err := suite.tool.readLiveVersion()
	suite.Require().NoError(err)
	suite.Equal(len(suite.tool.migrations)-1, v.Version)
}
