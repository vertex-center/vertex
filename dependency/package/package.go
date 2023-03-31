package _package

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/goccy/go-json"
	"github.com/vertex-center/vertex/storage"
)

const (
	PmNone   = "sources"
	PmAptGet = "apt-get"
	PmBrew   = "brew"
	PmPacman = "pacman"
	PmSnap   = "snap"
)

var pkgs map[string]Package

type Package struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Homepage       string            `json:"homepage"`
	License        string            `json:"license"`
	Check          string            `json:"check"`
	InstallPackage map[string]string `json:"install"`
}

func Read(id string) (*Package, error) {
	p := path.Join(GetPath(id), fmt.Sprintf("%s.json", id))

	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var pkg Package
	err = json.Unmarshal(file, &pkg)
	return &pkg, err
}

func GetPath(id string) string {
	return path.Join(storage.PathDependencies, "packages", id)
}

func (p *Package) Install(pm string) error {
	cmd := p.InstallationCommand(pm)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (p *Package) InstallationCommand(pm string) *exec.Cmd {
	if strings.HasPrefix(p.InstallPackage[pm], "script:") {
		script := strings.Split(p.InstallPackage[pm], ":")[1]
		return exec.Command(script)
	}

	packageName := p.InstallPackage[pm]

	switch pm {
	case PmAptGet:
		return exec.Command("apt-get", "install", packageName)
	case PmBrew:
		return exec.Command("brew", "install", packageName)
	case PmPacman:
		return exec.Command("pacman", "-S", "--noconfirm", packageName)
	case PmSnap:
		return exec.Command("snap", "install", packageName)
	}
	return nil
}

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

func Get(id string) (*Package, error) {
	pkg, ok := pkgs[id]
	if !ok {
		return nil, fmt.Errorf("the dependency %s was not found", id)
	}
	return &pkg, nil
}
