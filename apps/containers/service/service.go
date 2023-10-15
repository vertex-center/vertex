package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
)

type ServiceService struct {
	uuid           uuid.UUID
	serviceAdapter types.ServiceAdapterPort
}

func NewServiceService() *ServiceService {
	return &ServiceService{
		uuid:           uuid.New(),
		serviceAdapter: adapter.NewServiceFSAdapter(nil),
	}
}

func (s *ServiceService) GetById(id string) (types.Service, error) {
	return s.serviceAdapter.Get(id)
}

func (s *ServiceService) GetAll() []types.Service {
	return s.serviceAdapter.GetAll()
}

func (s *ServiceService) Reload() error {
	return s.serviceAdapter.Reload()
}

func (s *ServiceService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *ServiceService) OnEvent(e interface{}) {
	switch e.(type) {
	case vtypes.EventVertexUpdated:
		err := s.Reload()
		if err != nil {
			log.Error(err)
			return
		}
	}
}
