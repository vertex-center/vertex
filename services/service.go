package services

import (
	"github.com/vertex-center/vertex/repository"
	"github.com/vertex-center/vertex/types"
)

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
