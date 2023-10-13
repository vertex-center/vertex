package migration

import (
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type Migration1TestSuite struct {
	suite.Suite

	dir         string
	containerId string
}

func TestMigration1TestSuite(t *testing.T) {
	suite.Run(t, new(Migration1TestSuite))
}

func (suite *Migration1TestSuite) SetupSuite() {
	var err error
	suite.dir, err = os.MkdirTemp("", "migration_1_test-*")
	suite.Require().NoError(err)

	suite.containerId = uuid.New().String()
	err = os.MkdirAll(path.Join(suite.dir, "instances", suite.containerId), os.ModePerm)
	suite.Require().NoError(err)
}

func (suite *Migration1TestSuite) TearDownSuite() {
	err := os.RemoveAll("./migration_1_test-*")
	suite.NoError(err)
}

func (suite *Migration1TestSuite) TestUp() {
	m := &migration1{}
	err := m.Up(suite.dir)
	suite.NoError(err)

	_, err = os.Stat(path.Join(suite.dir, "instances"))
	suite.True(os.IsNotExist(err))

	_, err = os.Stat(path.Join(suite.dir, "apps", "vx-containers", suite.containerId))
	suite.NoError(err)
}

func (suite *Migration1TestSuite) TestUpNoLive() {
	m := &migration1{}

	dir, err := os.MkdirTemp("", "migration_1_test_empty-*")
	suite.Require().NoError(err)

	err = m.Up(dir)
	suite.NoError(err)
}
