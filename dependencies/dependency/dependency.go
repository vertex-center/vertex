package dependency

import (
	"fmt"
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

	if !strings.HasPrefix(pkg.Check, "script:") {
		return nil, fmt.Errorf("the value '%s' has no script: prefix", pkg.Check)
	}

	script := strings.Split(pkg.Check, ":")[1]

	cmd := exec.Command(path.Join(p, script))

	err = cmd.Run()
	if cmd.ProcessState.ExitCode() == 1 {
		return &Dependency{
			Package:   pkg,
			Installed: false,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &Dependency{
		Package:   pkg,
		Installed: true,
	}, nil
}
