package adapter

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrPkgNotFound = errors.New("package not found")
)

type PackageFSAdapter struct {
	pkgs      []types.Package
	pkgsMutex *sync.RWMutex

	dependenciesPath string
}

type PackageFSAdapterParams struct {
	dependenciesPath string
}

func NewPackageFSAdapter(params *PackageFSAdapterParams) types.PackageAdapterPort {
	if params == nil {
		params = &PackageFSAdapterParams{}
	}
	if params.dependenciesPath == "" {
		params.dependenciesPath = path.Join(storage.Path, "packages")
	}

	adapter := &PackageFSAdapter{
		pkgs:      []types.Package{},
		pkgsMutex: &sync.RWMutex{},

		dependenciesPath: params.dependenciesPath,
	}

	err := adapter.Reload()
	if err != nil {
		log.Error(fmt.Errorf("failed to reload services: %v", err))
	}
	return adapter
}

func (a *PackageFSAdapter) GetByID(id string) (types.Package, error) {
	a.pkgsMutex.RLock()
	defer a.pkgsMutex.RUnlock()

	for _, pkg := range a.pkgs {
		if pkg.ID == id {
			return pkg, nil
		}
	}

	return types.Package{}, ErrPkgNotFound
}

func (a *PackageFSAdapter) GetPath(id string) string {
	return path.Join(a.dependenciesPath, "packages", id)
}

func (a *PackageFSAdapter) Reload() error {
	dir, err := os.ReadDir(path.Join(a.dependenciesPath, "packages"))
	if err != nil {
		return err
	}

	a.pkgsMutex.RLock()
	defer a.pkgsMutex.RUnlock()

	a.pkgs = []types.Package{}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		pkg, err := a.readFromDisk(name)
		if err != nil {
			return err
		}

		a.pkgs = append(a.pkgs, *pkg)
	}

	return nil
}

func (a *PackageFSAdapter) readFromDisk(id string) (*types.Package, error) {
	p := path.Join(a.GetPath(id), fmt.Sprintf("%s.json", id))

	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var pkg types.Package
	err = json.Unmarshal(file, &pkg)
	pkg.ID = id
	return &pkg, err
}
