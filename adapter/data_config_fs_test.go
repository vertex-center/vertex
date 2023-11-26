package adapter

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DataConfigFSAdapterTestSuite struct {
	suite.Suite

	adapter *SettingsFSAdapter
}

func TestDataConfigFSAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(DataConfigFSAdapterTestSuite))
}

func (suite *DataConfigFSAdapterTestSuite) SetupTest() {
	suite.adapter = NewSettingsFSAdapter(&SettingsFSAdapterParams{
		settingsDir: suite.T().TempDir(),
	}).(*SettingsFSAdapter)
}

func (suite *DataConfigFSAdapterTestSuite) TestReadDataConfig() {
	data, err := json.Marshal(suite.adapter.settings)
	suite.Require().NoError(err)

	p := path.Join(suite.adapter.settingsDir, "data_config.yml")
	err = os.WriteFile(p, data, 0644)
	suite.Require().NoError(err)

	err = suite.adapter.read()
	suite.Require().NoError(err)
}
