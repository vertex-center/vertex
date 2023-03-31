package _package

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	PathLive         = "live_test"
	PathDependencies = "live_test/dependencies"
)

type PackageTestSuite struct {
	suite.Suite
}

func TestPackageTestSuite(t *testing.T) {
	suite.Run(t, new(PackageTestSuite))
}

func (suite *PackageTestSuite) SetupSuite() {
	testReload(suite.T())
}

func (suite *PackageTestSuite) TearDownSuite() {
	err := os.RemoveAll(PathLive)
	assert.NoError(suite.T(), err)
}

func testReload(t *testing.T) {
	err := os.MkdirAll(PathDependencies, os.ModePerm)
	assert.NoError(t, err)

	// reload to test Clone()
	err = reload(PathDependencies)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(pkgs))
	assert.Equal(t, "Redis", pkgs["redis"].Name)
	assert.Equal(t, "BSD-3", pkgs["redis"].License)

	// reload to test Pull()
	err = reload(PathDependencies)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(pkgs))
	assert.Equal(t, "Redis", pkgs["redis"].Name)
	assert.Equal(t, "BSD-3", pkgs["redis"].License)
}

func (suite *PackageTestSuite) TestGetPath() {
	p := getPath(PathDependencies, "redis")
	assert.Equal(suite.T(), "live_test/dependencies/packages/redis", p)
}

func (suite *PackageTestSuite) TestReload() {
	err := reload(PathDependencies)
	assert.NoError(suite.T(), err)
}

func (suite *PackageTestSuite) TestGet() {
	pkg, err := Get("redis")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.Name, "Redis")

	pkg, err = Get("undefined_package_name")
	assert.ErrorIs(suite.T(), err, ErrPkgNotFound)
}
