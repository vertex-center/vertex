package migration

import (
	"os"
	"path"
)

// migration1 renames the 'instances' directory to 'apps/vx-containers'.
type migration1 struct{}

func (m *migration1) Up(livePath string) error {
	instancesPath := path.Join(livePath, "instances")
	containersPath := path.Join(livePath, "apps", "vx-containers")

	err := os.MkdirAll(path.Dir(containersPath), 0755)
	if err != nil {
		return err
	}

	err = os.Rename(instancesPath, containersPath)
	if err != nil && os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func (m *migration1) DispatchCommands() []interface{} {
	return []interface{}{
		CommandRecreateContainers{},
	}
}
