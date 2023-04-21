package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex/repository"
	"github.com/vertex-center/vertex/types"
)

var logger = console.New("vertex::services")

type ServiceService struct {
	repo repository.ServiceRepository
}

func NewServiceService() ServiceService {
	return ServiceService{
		repo: repository.NewServiceRepository(nil),
	}
}

func (s *ServiceService) ListAvailable() []types.Service {
	return s.repo.GetAvailable()
}

func ReadFromDisk(servicePath string) (*types.Service, error) {
	// TODO: Move this method elsewhere

	data, err := os.ReadFile(path.Join(servicePath, ".vertex", "service.json"))
	if err != nil {
		logger.Warn(fmt.Sprintf("service at '%s' has no '.vertex/service.json' file", path.Dir(servicePath)))
	}

	var service types.Service
	err = json.Unmarshal(data, &service)
	return &service, err
}
