package updates

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vertex-center/vertex/core/types"
)

type VertexUpdaterTestSuite struct {
	suite.Suite
	updater VertexUpdater
}

func TestVertexUpdaterTestSuite(t *testing.T) {
	suite.Run(t, new(VertexUpdaterTestSuite))
}

func (suite *VertexUpdaterTestSuite) SetupTest() {
	suite.updater = NewVertexUpdater(types.About{
		Version: "0.12.0",
	})
}

func (suite *VertexUpdaterTestSuite) TestCurrentVersion() {
	version, err := suite.updater.CurrentVersion()
	suite.NoError(err)
	suite.Equal("v0.12.0", version)
}

func (suite *VertexUpdaterTestSuite) TestID() {
	suite.Equal("vertex", suite.updater.ID())
}
