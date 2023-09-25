package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	PathPackages = "tests/packages"
)

type PackageAdapterTestSuite struct {
	suite.Suite

	adapter PackageFSAdapter
}

func TestPackageAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(PackageAdapterTestSuite))
}

func (suite *PackageAdapterTestSuite) SetupSuite() {
	suite.adapter = *NewPackageFSAdapter(&PackageFSAdapterParams{
		dependenciesPath: PathPackages,
	}).(*PackageFSAdapter)

	err := suite.adapter.Reload()
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), 0, len(suite.adapter.pkgs))

	pkg, err := suite.adapter.GetByID("redis")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "redis", pkg.ID)
	assert.Equal(suite.T(), "Redis", pkg.Name)
	assert.Equal(suite.T(), "BSD-3", pkg.License)
}

func (suite *PackageAdapterTestSuite) TestGetPath() {
	p := suite.adapter.GetPath("redis")
	assert.Equal(suite.T(), "tests/packages/packages/redis", p)
}

func (suite *PackageAdapterTestSuite) TestReload() {
	err := suite.adapter.Reload()
	assert.NoError(suite.T(), err)
}

func (suite *PackageAdapterTestSuite) TestGet() {
	pkg, err := suite.adapter.GetByID("redis")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.Name, "Redis")

	_, err = suite.adapter.GetByID("undefined_package_name")
	assert.ErrorIs(suite.T(), err, ErrPkgNotFound)
}
