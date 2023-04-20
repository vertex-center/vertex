package dependency

import (
	"os/exec"
	"path"
	"strings"

	"github.com/vertex-center/vertex/services/pkg"
	"github.com/vertex-center/vertex/types"
)

func Get(id string) (*types.Dependency, error) {
	p, err := pkg.Get(id)
	if err != nil {
		return nil, err
	}

	pkgPath := pkg.GetPath(id)

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
			return nil, err
		}
	} else {
		_, err := exec.LookPath(p.Check)
		if err == nil {
			installed = true
		}
	}

	return &types.Dependency{
		Package:   p,
		Installed: installed,
	}, nil
}
