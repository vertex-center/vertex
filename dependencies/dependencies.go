package dependencies

import (
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/vertex-center/vertex/storage"
)

func Reload() error {
	_, err := git.PlainClone(storage.PathDependencies, false, &git.CloneOptions{
		URL:      "https://github.com/vertex-center/vertex-dependencies",
		Progress: os.Stdout,
	})

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		repo, err := git.PlainOpen(storage.PathDependencies)
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
