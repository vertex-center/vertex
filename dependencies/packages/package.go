package packages

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

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
		return exec.Command("pacman", "-S", packageName)
	case PmSnap:
		return exec.Command("snap", "install", packageName)
	}
	return nil
}
