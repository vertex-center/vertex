package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrPkgNotFound = errors.New("package not found")
)

type PackageFSRepository struct {
	pkgs      map[string]types.Package
	pkgsMutex *sync.RWMutex

	dependenciesPath string
}

type PackageRepositoryParams struct {
	dependenciesPath string
}

func NewPackageFSRepository(params *PackageRepositoryParams) PackageFSRepository {
	if params == nil {
		params = &PackageRepositoryParams{}
	}
	if params.dependenciesPath == "" {
		params.dependenciesPath = path.Join(storage.Path, "packages")
	}

	repo := PackageFSRepository{
		pkgs:      map[string]types.Package{},
		pkgsMutex: &sync.RWMutex{},

		dependenciesPath: params.dependenciesPath,
	}

	err := repo.Reload()
	if err != nil {
		logger.Error(fmt.Errorf("failed to reload services repository: %v", err)).Print()
	}
	return repo
}

func (r *PackageFSRepository) GetByID(id string) (types.Package, error) {
	r.pkgsMutex.RLock()
	defer r.pkgsMutex.RUnlock()

	pkg, ok := r.pkgs[id]
	if !ok {
		return types.Package{}, ErrPkgNotFound
	}
	return pkg, nil
}

func (r *PackageFSRepository) set(id string, pkg types.Package) {
	r.pkgsMutex.Lock()
	defer r.pkgsMutex.Unlock()

	r.pkgs[id] = pkg
}

func (r *PackageFSRepository) GetPath(id string) string {
	return path.Join(r.dependenciesPath, "packages", id)
}

func (r *PackageFSRepository) Reload() error {
	dir, err := os.ReadDir(path.Join(r.dependenciesPath, "packages"))
	if err != nil {
		return err
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		pkg, err := r.readFromDisk(name)
		if err != nil {
			return err
		}

		r.set(name, *pkg)
	}

	return nil
}

func (r *PackageFSRepository) readFromDisk(id string) (*types.Package, error) {
	p := path.Join(r.GetPath(id), fmt.Sprintf("%s.json", id))

	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var pkg types.Package
	err = json.Unmarshal(file, &pkg)
	return &pkg, err
}
