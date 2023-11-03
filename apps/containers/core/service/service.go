package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	evtypes "github.com/vertex-center/vertex/pkg/event/types"
	"github.com/vertex-center/vertex/pkg/log"
)

type ServiceService struct {
	uuid           uuid.UUID
	serviceAdapter port.ServiceAdapter
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

func (s *ServiceService) reload() error {
	return s.serviceAdapter.Reload()
}

func (s *ServiceService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *ServiceService) OnEvent(e evtypes.Event) {
	switch e.(type) {
	case vtypes.EventVertexUpdated:
		err := s.reload()
		if err != nil {
			log.Error(err)
			return
		}
	}
}
