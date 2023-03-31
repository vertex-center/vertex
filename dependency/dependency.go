package dependency

import (
	"os/exec"
	"path"
	"strings"

	"github.com/vertex-center/vertex/dependency/package"
)

type Dependency struct {
	*_package.Package

	Installed bool `json:"installed"`
}

func Get(id string) (*Dependency, error) {
	pkg, err := _package.Get(id)
	if err != nil {
		return nil, err
	}

	p := _package.GetPath(id)

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
