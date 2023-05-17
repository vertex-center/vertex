package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	errors2 "github.com/pkg/errors"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
)

var (
	ErrPkgNotFound = errors2.New("package not found")
)

type PackageFSRepository struct {
	pkgs             map[string]types.Package
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
		params.dependenciesPath = storage.PathPackages
	}

	repo := PackageFSRepository{
		dependenciesPath: params.dependenciesPath,
		pkgs:             map[string]types.Package{},
	}
	err := repo.reload()
	if err != nil {
		logger.Error(fmt.Errorf("failed to reload services repository: %v", err)).Print()
	}
	return repo
}

func (r *PackageFSRepository) Get(id string) (types.Package, error) {
	pkg, ok := r.pkgs[id]
	if !ok {
		return types.Package{}, ErrPkgNotFound
	}
	return pkg, nil
}

func (r *PackageFSRepository) GetPath(id string) string {
	return path.Join(r.dependenciesPath, "packages", id)
}

func (r *PackageFSRepository) reload() error {
	dir, err := os.ReadDir(path.Join(r.dependenciesPath, "packages"))
	if err != nil {
		return err
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		pkg, err := r.readPkgFromDisk(name)
		if err != nil {
			return err
		}

		r.pkgs[name] = *pkg
	}

	return nil
}

func (r *PackageFSRepository) readPkgFromDisk(id string) (*types.Package, error) {
	p := path.Join(r.GetPath(id), fmt.Sprintf("%s.json", id))

	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var pkg types.Package
	err = json.Unmarshal(file, &pkg)
	return &pkg, err
}
