package updates

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type VertexClientUpdaterTestSuite struct {
	suite.Suite
	updater VertexClientUpdater
}

func TestVertexClientUpdaterTestSuite(t *testing.T) {
	suite.Run(t, new(VertexClientUpdaterTestSuite))
}

func (suite *VertexClientUpdaterTestSuite) SetupTest() {
	temp := suite.T().TempDir()

	suite.updater = NewVertexClientUpdater(temp)

	err := os.MkdirAll(path.Join(temp, "dist"), os.ModePerm)
	suite.Require().NoError(err)

	err = os.WriteFile(path.Join(temp, "dist", "version.txt"), []byte("v0.12.0"), os.ModePerm)
	suite.Require().NoError(err)
}

func (suite *VertexClientUpdaterTestSuite) TestCurrentVersion() {
	version, err := suite.updater.CurrentVersion()
	suite.Require().NoError(err)
	suite.Equal("v0.12.0", version)
}

func (suite *VertexClientUpdaterTestSuite) TestID() {
	suite.Equal("vertex_client", suite.updater.ID())
}
