package repository

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"gopkg.in/yaml.v2"
)

type ServiceFSRepository struct {
	servicesPath string
	services     []types.Service
}

type ServiceRepositoryParams struct {
	servicesPath string
}

func NewServiceFSRepository(params *ServiceRepositoryParams) ServiceFSRepository {
	if params == nil {
		params = &ServiceRepositoryParams{}
	}
	if params.servicesPath == "" {
		params.servicesPath = path.Join(storage.Path, "services")
	}

	repo := ServiceFSRepository{
		servicesPath: params.servicesPath,
	}
	err := repo.Reload()
	if err != nil {
		log.Default.Error(fmt.Errorf("failed to reload services repository: %v", err))
	}
	return repo
}

func (r *ServiceFSRepository) Get(id string) (types.Service, error) {
	for _, service := range r.services {
		if service.ID == id {
			return service, nil
		}
	}

	return types.Service{}, types.ErrServiceNotFound
}

func (r *ServiceFSRepository) GetScript(id string) ([]byte, error) {
	service, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	if service.Methods.Script == nil {
		return nil, errors.New("the service doesn't have a script method")
	}

	return os.ReadFile(path.Join(r.servicesPath, "services", id, service.Methods.Script.Filename))
}

func (r *ServiceFSRepository) GetAll() []types.Service {
	return r.services
}

func (r *ServiceFSRepository) Reload() error {
	servicesPath := path.Join(r.servicesPath, "services")

	r.services = []types.Service{}

	entries, err := os.ReadDir(servicesPath)
	if err != nil {
		return err
	}

	for _, dir := range entries {
		if !dir.IsDir() {
			continue
		}

		servicePath := path.Join(servicesPath, dir.Name(), "service.yml")

		file, err := os.ReadFile(servicePath)
		if err != nil {
			return err
		}

		var service types.Service
		err = yaml.Unmarshal(file, &service)
		if err != nil {
			return err
		}

		r.services = append(r.services, service)
	}

	return nil
}
