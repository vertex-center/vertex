package migration

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type Migration0TestSuite struct {
	suite.Suite

	dir        string
	instanceId string
}

func TestMigration0TestSuite(t *testing.T) {
	suite.Run(t, new(Migration0TestSuite))
}

func (suite *Migration0TestSuite) SetupSuite() {
	var err error
	suite.dir, err = os.MkdirTemp("", "migration_0_test-*")
	suite.Require().NoError(err)

	suite.instanceId = uuid.New().String()
	p := path.Join(suite.dir, "instances", suite.instanceId, ".vertex")
	err = os.MkdirAll(p, os.ModePerm)
	suite.Require().NoError(err)

	j, err := json.Marshal(map[string]string{"a": "b"})
	suite.Require().NoError(err)

	err = os.WriteFile(path.Join(p, "instance_settings.json"), j, os.ModePerm)
	suite.Require().NoError(err)
}

func (suite *Migration0TestSuite) TearDownSuite() {
	err := os.RemoveAll("./migration_0_test-*")
	suite.Require().NoError(err)
}

func (suite *Migration0TestSuite) TestUp() {
	m := &migration0{}
	err := m.Up(suite.dir)
	suite.Require().NoError(err)

	_, err = os.Stat(path.Join(suite.dir, "instances", suite.instanceId, ".vertex", "instance_settings.json"))
	suite.True(os.IsNotExist(err))

	_, err = os.Stat(path.Join(suite.dir, "instances", suite.instanceId, ".vertex", "settings.yml"))
	suite.Require().NoError(err)
}

func (suite *Migration0TestSuite) TestUpNoLive() {
	m := &migration0{}

	dir, err := os.MkdirTemp("", "migration_0_test_empty-*")
	suite.Require().NoError(err)

	err = m.Up(dir)
	suite.Require().NoError(err)
}
