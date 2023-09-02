package services

import (
	"os"
	"os/exec"
	"path"
	"strings"

	errors2 "github.com/pkg/errors"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrPkgManagerNotFound = errors2.New("package manager not found")
)

type PackageService struct {
	packageRepo types.PackageRepository
}

func NewPackageService(packageRepo types.PackageRepository) PackageService {
	return PackageService{
		packageRepo: packageRepo,
	}
}

type InstallCmd struct {
	Cmd  string
	Sudo bool
}

func (s *PackageService) Install(c InstallCmd) error {
	args := strings.Fields(c.Cmd)
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (s *PackageService) InstallationCommand(p *types.Package, pm string) (InstallCmd, error) {
	if strings.HasPrefix(p.InstallPackage[pm], "script:") {
		return InstallCmd{
			Cmd:  strings.Split(p.InstallPackage[pm], ":")[1],
			Sudo: false,
		}, nil
	}

	packageName := p.InstallPackage[pm]

	switch pm {
	case types.PmAptGet:
		return InstallCmd{
			Cmd:  "sudo apt-get install " + packageName,
			Sudo: true,
		}, nil
	case types.PmBrew:
		return InstallCmd{
			Cmd:  "brew install " + packageName,
			Sudo: false,
		}, nil
	case types.PmNpm:
		return InstallCmd{
			Cmd:  "npm install -g " + packageName,
			Sudo: false,
		}, nil
	case types.PmPacman:
		return InstallCmd{
			Cmd:  "sudo pacman -S --noconfirm " + packageName,
			Sudo: true,
		}, nil
	case types.PmSnap:
		return InstallCmd{
			Cmd:  "sudo snap install " + packageName,
			Sudo: true,
		}, nil
	}

	return InstallCmd{}, ErrPkgManagerNotFound
}

func (s *PackageService) Get(id string) (types.Package, error) {
	p, err := s.packageRepo.Get(id)
	if err != nil {
		return types.Package{}, err
	}

	pkgPath := s.packageRepo.GetPath(id)

	isScript := strings.HasPrefix(p.Check, "script:")
	installed := false

	if isScript {
		script := strings.Split(p.Check, ":")[1]

		cmd := exec.Command(path.Join(pkgPath, script))

		err = cmd.Run()
		if cmd.ProcessState.ExitCode() == 0 {
			installed = true
		}
		if err != nil {
			return types.Package{}, err
		}
	} else {
		_, err := exec.LookPath(p.Check)
		if err == nil {
			installed = true
		}
	}

	p.Installed = &installed
	return p, nil
}

func (s *PackageService) Reload() error {
	return s.packageRepo.Reload()
}
