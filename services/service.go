package services

import (
	"github.com/vertex-center/vertex/types"
)

type ServiceService struct {
	serviceAdapter types.ServiceAdapterPort
}

func NewServiceService(serviceAdapter types.ServiceAdapterPort) ServiceService {
	return ServiceService{
		serviceAdapter: serviceAdapter,
	}
}

func (s *ServiceService) ListAvailable() []types.Service {
	return s.serviceAdapter.GetAll()
}

func (s *ServiceService) Reload() error {
	return s.serviceAdapter.Reload()
}
