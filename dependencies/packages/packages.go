package packages

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/vertex-center/vertex/storage"
)

var pkgs map[string]Package

func Reload() error {
	pkgs = map[string]Package{}

	err := ReloadRepository()
	if err != nil {
		return err
	}

	dir, err := os.ReadDir(path.Join(storage.PathDependencies, "packages"))
	if err != nil {
		return err
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		pkg, err := Read(name)
		if err != nil {
			return err
		}

		pkgs[name] = *pkg
	}

	return nil
}

func Get(id string) (*Package, error) {
	pkg, ok := pkgs[id]
	if !ok {
		return nil, fmt.Errorf("the dependency %s was not found", id)
	}
	return &pkg, nil
}

func ReloadRepository() error {
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
