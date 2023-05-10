package repository

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/vertex-center/vertex/storage"
	"github.com/vertex-center/vertex/types"
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
		params.servicesPath = storage.PathServices
	}

	repo := ServiceFSRepository{
		servicesPath: params.servicesPath,
	}
	err := repo.reload()
	if err != nil {
		log.Fatalf("failed to reload services repository: %v", err)
	}
	return repo
}

func (r *ServiceFSRepository) GetAll() []types.Service {
	return r.services
}

func (r *ServiceFSRepository) reload() error {
	url := "https://github.com/vertex-center/vertex-services"

	err := storage.CloneOrPullRepository(url, r.servicesPath)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(path.Join(r.servicesPath, "services.json"))
	if err != nil {
		return err
	}

	var availableMap map[string]types.Service
	err = json.Unmarshal(file, &availableMap)
	if err != nil {
		return err
	}

	r.services = []types.Service{}
	for key, service := range availableMap {
		service.ID = key
		r.services = append(r.services, service)
	}

	return nil
}
