package storage

import (
	"os"
	"testing"

	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestCloneRepository() {
	fs := fixtures.Basic().One().DotGit()

	dir, err := os.MkdirTemp("", "*_live_test")
	suite.NoError(err)

	defer os.RemoveAll(dir)

	err = CloneRepository(fs.Root(), dir)
	suite.NoError(err)
	suite.DirExists(dir)

	err = CloneRepository(fs.Root(), dir)
	suite.ErrorIs(err, git.ErrRepositoryAlreadyExists)
}

func (suite *RepositoryTestSuite) TestCloneOrPullRepository() {
	fs := fixtures.Basic().One().DotGit()

	dir, err := os.MkdirTemp("", "*_live_test")
	suite.NoError(err)

	defer os.RemoveAll(dir)

	// reload to test Clone()
	err = CloneOrPullRepository(fs.Root(), dir)
	suite.NoError(err)
	suite.DirExists(dir)

	// reload to test Pull()
	err = CloneOrPullRepository(fs.Root(), dir)
	suite.NoError(err)
	suite.DirExists(dir)
}
