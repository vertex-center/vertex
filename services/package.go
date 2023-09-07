package services

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrPkgManagerNotFound = errors.New("package manager not found")
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

func (s *PackageService) GetByID(id string) (types.Package, error) {
	pkg, err := s.packageRepo.GetByID(id)
	if err != nil {
		return types.Package{}, err
	}

	installed, err := s.checkIsInstalled(id, pkg)
	pkg.Installed = &installed
	return pkg, err
}

func (s *PackageService) checkIsInstalled(id string, pkg types.Package) (bool, error) {
	pkgPath := s.packageRepo.GetPath(id)
	isScript := strings.HasPrefix(pkg.Check, "script:")

	if isScript {
		return s.checkIsInstalledWithScript(pkgPath, pkg.Check)
	}
	return s.checkIsInstalledWithCommand(pkg.Check)
}

func (s *PackageService) checkIsInstalledWithScript(pkgPath, check string) (bool, error) {
	script := strings.Split(check, ":")[1]
	cmd := exec.Command(path.Join(pkgPath, script))

	err := cmd.Run()
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		return false, err
	}
	return true, nil
}

func (s *PackageService) checkIsInstalledWithCommand(cmd string) (bool, error) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *PackageService) Reload() error {
	return s.packageRepo.Reload()
}
