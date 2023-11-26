package adapter

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type DataConfigFSAdapterTestSuite struct {
	suite.Suite

	adapter *DataConfigFSAdapter
}

func TestDataConfigFSAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(DataConfigFSAdapterTestSuite))
}

func (suite *DataConfigFSAdapterTestSuite) SetupTest() {
	suite.adapter = NewDataConfigFSAdapter(&DataConfigFSAdapterParams{
		configDir: suite.T().TempDir(),
	}).(*DataConfigFSAdapter)
}

func (suite *DataConfigFSAdapterTestSuite) TestReadDataConfig() {
	data, err := yaml.Marshal(suite.adapter.config)
	suite.Require().NoError(err)

	p := path.Join(suite.adapter.configDir, "config.yml")
	err = os.WriteFile(p, data, 0644)
	suite.Require().NoError(err)

	err = suite.adapter.read()
	suite.Require().NoError(err)
}
