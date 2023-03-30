package dependency

import (
	"os/exec"
	"path"
	"strings"

	"github.com/vertex-center/vertex/dependencies/packages"
)

type Dependency struct {
	*packages.Package

	Installed bool `json:"installed"`
}

func New(id string) (*Dependency, error) {
	pkg, err := packages.Get(id)
	if err != nil {
		return nil, err
	}

	p := packages.GetPath(id)

	isScript := strings.HasPrefix(pkg.Check, "script:")
	installed := false

	if isScript {
		script := strings.Split(pkg.Check, ":")[1]

		cmd := exec.Command(path.Join(p, script))

		err = cmd.Run()
		if cmd.ProcessState.ExitCode() == 0 {
			installed = true
		}
		if err != nil {
			return nil, err
		}
	} else {
		_, err := exec.LookPath(pkg.Check)
		if err == nil {
			installed = true
		}
	}

	return &Dependency{
		Package:   pkg,
		Installed: installed,
	}, nil
}
