package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/goccy/go-json"
	errors2 "github.com/pkg/errors"
	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
)

const (
	PmNone   = "sources"
	PmAptGet = "apt-get"
	PmBrew   = "brew"
	PmPacman = "pacman"
	PmSnap   = "snap"
)

var (
	ErrPkgNotFound        = errors2.New("package not found")
	ErrPkgManagerNotFound = errors2.New("package manager not found")
)

var pkgs map[string]types.Package

type InstallCmd struct {
	Cmd  string
	Sudo bool
}

func (c *InstallCmd) Exec() error {
	args := strings.Fields(c.Cmd)
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func read(pathDependencies string, id string) (*types.Package, error) {
	p := path.Join(getPath(pathDependencies, id), fmt.Sprintf("%s.json", id))

	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var pkg types.Package
	err = json.Unmarshal(file, &pkg)
	return &pkg, err
}

func GetPath(packageID string) string {
	return getPath(storage.PathDependencies, packageID)
}

func getPath(pathDependencies string, id string) string {
	return path.Join(pathDependencies, "packages", id)
}

func InstallationCommand(p *types.Package, pm string) (InstallCmd, error) {
	if strings.HasPrefix(p.InstallPackage[pm], "script:") {
		return InstallCmd{
			Cmd:  strings.Split(p.InstallPackage[pm], ":")[1],
			Sudo: false,
		}, nil
	}

	packageName := p.InstallPackage[pm]

	switch pm {
	case PmAptGet:
		return InstallCmd{
			Cmd:  "sudo apt-get install " + packageName,
			Sudo: true,
		}, nil
	case PmBrew:
		return InstallCmd{
			Cmd:  "brew install " + packageName,
			Sudo: false,
		}, nil
	case PmPacman:
		return InstallCmd{
			Cmd:  "sudo pacman -S --noconfirm " + packageName,
			Sudo: true,
		}, nil
	case PmSnap:
		return InstallCmd{
			Cmd:  "sudo snap install " + packageName,
			Sudo: true,
		}, nil
	}

	return InstallCmd{}, ErrPkgManagerNotFound
}

func Reload() error {
	return reload(storage.PathDependencies)
}

func reload(dependenciesPath string) error {
	pkgs = map[string]types.Package{}

	url := "https://github.com/vertex-center/vertex-dependencies"

	err := storage.CloneOrPullRepository(url, dependenciesPath)
	if err != nil {
		return err
	}

	dir, err := os.ReadDir(path.Join(dependenciesPath, "packages"))
	if err != nil {
		return err
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		pkg, err := read(dependenciesPath, name)
		if err != nil {
			return err
		}

		pkgs[name] = *pkg
	}

	return nil
}

func Get(id string) (*types.Package, error) {
	pkg, ok := pkgs[id]
	if !ok {
		return nil, ErrPkgNotFound
	}
	return &pkg, nil
}
