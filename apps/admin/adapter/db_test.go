package adapter

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type DbConfigFSAdapterTestSuite struct {
	suite.Suite

	adapter *DbAdapter
}

func TestDataConfigFSAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(DbConfigFSAdapterTestSuite))
}

func (suite *DbConfigFSAdapterTestSuite) SetupTest() {
	suite.adapter = NewDbAdapter(&DbAdapterParams{
		configDir: suite.T().TempDir(),
	}).(*DbAdapter)
}

func (suite *DbConfigFSAdapterTestSuite) TestReadDataConfig() {
	data, err := yaml.Marshal(suite.adapter.config)
	suite.Require().NoError(err)

	p := path.Join(suite.adapter.configDir, "config.yml")
	err = os.WriteFile(p, data, 0644)
	suite.Require().NoError(err)

	err = suite.adapter.read()
	suite.Require().NoError(err)
}
