package services

import (
	"github.com/vertex-center/vertex/types"
)

type ServiceService struct {
	serviceRepo types.ServiceRepository
}

func NewServiceService(serviceRepo types.ServiceRepository) ServiceService {
	return ServiceService{
		serviceRepo: serviceRepo,
	}
}

func (s *ServiceService) ListAvailable() []types.Service {
	return s.serviceRepo.GetAll()
}
