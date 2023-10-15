package updates

import (
	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositoryUpdaterTestSuite struct {
	suite.Suite
	updater RepositoryUpdater
	repo    *git.Repository
}

func TestRepositoryUpdaterTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryUpdaterTestSuite))
}

func (suite *RepositoryUpdaterTestSuite) SetupTest() {
	dir := suite.T().TempDir()

	suite.updater = NewRepositoryUpdater("id", dir, "owner", "repo")

	fs := fixtures.Basic().One().DotGit()

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: fs.Root(),
	})
	suite.Require().NoError(err)

	suite.repo, err = git.PlainOpen(dir)
	suite.Require().NoError(err)
}

func (suite *RepositoryUpdaterTestSuite) TestCurrentVersion() {
	worktree, err := suite.repo.Worktree()
	suite.Require().NoError(err)

	hash, err := worktree.Commit("test", &git.CommitOptions{
		All:               true,
		AllowEmptyCommits: true,
		Author: &object.Signature{
			Name:  "test",
			Email: "test@test.test",
		},
	})
	suite.NoError(err)

	version, err := suite.updater.CurrentVersion()
	suite.NoError(err)
	suite.Equal(hash.String(), version)
}

func (suite *RepositoryUpdaterTestSuite) TestID() {
	suite.Equal("id", suite.updater.ID())
}
