package services

import (
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
