package adapter

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SettingsFSAdapterTestSuite struct {
	suite.Suite

	adapter *SettingsFSAdapter
}

func TestSettingsFSAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsFSAdapterTestSuite))
}

func (suite *SettingsFSAdapterTestSuite) SetupTest() {
	suite.adapter = NewSettingsFSAdapter(&SettingsFSAdapterParams{
		settingsDir: suite.T().TempDir(),
	}).(*SettingsFSAdapter)
}

func (suite *SettingsFSAdapterTestSuite) TestReadSettings() {
	data, err := json.Marshal(suite.adapter.settings)
	suite.Require().NoError(err)

	p := path.Join(suite.adapter.settingsDir, "settings.json")
	err = os.WriteFile(p, data, 0644)
	suite.Require().NoError(err)

	err = suite.adapter.read()
	suite.Require().NoError(err)
}

func (suite *SettingsFSAdapterTestSuite) TestReadNonExistingSettings() {
	err := suite.adapter.read()
	suite.Require().ErrorIs(err, errSettingsNotFound)
}

func (suite *SettingsFSAdapterTestSuite) TestReadCorruptedSettings() {
	p := path.Join(suite.adapter.settingsDir, "settings.json")
	data := []byte("{{{corrupted:json}")
	err := os.WriteFile(p, data, 0644)
	suite.Require().NoError(err)

	err = suite.adapter.read()
	suite.Require().ErrorIs(err, errSettingsFailedToDecode)
}
