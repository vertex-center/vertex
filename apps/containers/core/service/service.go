package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/common/event"
	"github.com/vertex-center/vertex/pkg/event"
)

type serviceService struct {
	uuid           uuid.UUID
	serviceAdapter port.ServiceAdapter
}

func NewServiceService() port.ServiceService {
	return &serviceService{
		uuid:           uuid.New(),
		serviceAdapter: adapter.NewServiceFSAdapter(nil),
	}
}

func (s *serviceService) GetById(id string) (types.Service, error) {
	return s.serviceAdapter.Get(id)
}

func (s *serviceService) GetAll() []types.Service {
	return s.serviceAdapter.GetAll()
}

func (s *serviceService) reload() error {
	return s.serviceAdapter.Reload()
}

func (s *serviceService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *serviceService) OnEvent(e event.Event) error {
	switch e.(type) {
	case vtypes.VertexUpdated:
		return s.reload()
	}
	return nil
}
