package updates

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
	"os"
	"path"
)

type RepositoryUpdater struct {
	id    string
	dir   string
	owner string
	repo  string
}

func NewRepositoryUpdater(id, dir, owner, repo string) RepositoryUpdater {
	return RepositoryUpdater{
		id:    id,
		dir:   dir,
		owner: owner,
		repo:  repo,
	}
}

func (u RepositoryUpdater) CurrentVersion() (string, error) {
	repo, err := git.PlainOpen(u.dir)
	if err != nil {
		return "", nil
	}
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}

func (u RepositoryUpdater) Install(version string) error {
	url := fmt.Sprintf("https://github.com/%s/%s", u.owner, u.repo)

	log.Info("installing package", vlog.String("url", url), vlog.String("version", version))

	err := storage.CloneRepository(url, u.dir)
	if err != nil && !errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return err
	} else if err == nil {
		// This is freshly cloned, so we don't need to pull.
		return nil
	}

	repo, err := git.PlainOpen(u.dir)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	checkoutOptions := git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"),
		Force:  true,
	}

	err = worktree.Checkout(&checkoutOptions)
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(version),
		Force: true,
	})
}

func (u RepositoryUpdater) IsInstalled() bool {
	_, err := os.Stat(path.Join(u.dir))
	return err == nil
}

func (u RepositoryUpdater) ID() string {
	return u.id
}
