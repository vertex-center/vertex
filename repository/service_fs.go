package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
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
		logger.Error(fmt.Errorf("failed to reload services repository: %v", err)).Print()
	}
	return repo
}

func (r *ServiceFSRepository) Get(repo string) (types.Service, error) {
	if strings.HasPrefix(repo, "github.com") {
		repo = strings.TrimPrefix(repo, "github.com/")
	} else {
		return types.Service{}, errors.New("this repo is not supported")
	}

	println(repo)

	res, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/main/.vertex/service.json", repo))
	if err != nil {
		return types.Service{}, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return types.Service{}, err
	}

	var service types.Service
	err = json.Unmarshal(content, &service)
	if err != nil {
		return types.Service{}, err
	}

	return service, nil
}

func (r *ServiceFSRepository) GetAll() []types.Service {
	return r.services
}

func (r *ServiceFSRepository) reload() error {
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
