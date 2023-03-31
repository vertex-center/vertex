package storage

import (
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
)

func DownloadLatestRepository(repoPath string, url string) error {
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			return err
		}

		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = worktree.Pull(&git.PullOptions{})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
