package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	PathPackages = "package_test_live"
)

type PackageRepositoryTestSuite struct {
	suite.Suite

	repo PackageRepository
}

func TestPackageRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PackageRepositoryTestSuite))
}

func (suite *PackageRepositoryTestSuite) SetupSuite() {
	suite.repo = NewPackageRepo(&PackageRepositoryParams{
		dependenciesPath: PathPackages,
	})

	err := os.MkdirAll(PathPackages, os.ModePerm)
	assert.NoError(suite.T(), err)

	err = suite.repo.reload()
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), 0, len(suite.repo.pkgs))
	assert.Equal(suite.T(), "Redis", suite.repo.pkgs["redis"].Name)
	assert.Equal(suite.T(), "BSD-3", suite.repo.pkgs["redis"].License)
}

func (suite *PackageRepositoryTestSuite) TearDownSuite() {
	err := os.RemoveAll(PathPackages)
	assert.NoError(suite.T(), err)
}

func (suite *PackageRepositoryTestSuite) TestGetPath() {
	p := suite.repo.GetPkgPath("redis")
	assert.Equal(suite.T(), "package_test_live/packages/redis", p)
}

func (suite *PackageRepositoryTestSuite) TestReload() {
	err := suite.repo.reload()
	assert.NoError(suite.T(), err)
}

func (suite *PackageRepositoryTestSuite) TestGet() {
	pkg, err := suite.repo.Get("redis")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.Name, "Redis")

	_, err = suite.repo.Get("undefined_package_name")
	assert.ErrorIs(suite.T(), err, ErrPkgNotFound)
}
