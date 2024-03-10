package service

import (
	"context"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
)

type portsService struct {
	ports port.PortAdapter
}

func NewPortsService(ports port.PortAdapter) port.PortsService {
	return &portsService{ports}
}

func (s *portsService) GetPorts(ctx context.Context, filters types.PortFilters) (types.Ports, error) {
	return s.ports.GetPorts(ctx, filters)
}

func (s *portsService) PatchPort(ctx context.Context, p types.Port) error {
	err := p.Validate()
	if err != nil {
		return err
	}
	return s.ports.UpdatePortByID(ctx, p)
}

func (s *portsService) DeletePort(ctx context.Context, id uuid.UUID) error {
	return s.ports.DeletePort(ctx, id)
}

func (s *portsService) CreatePort(ctx context.Context, p types.Port) error {
	p.ID = uuid.New()
	err := p.Validate()
	if err != nil {
		return err
	}
	return s.ports.CreatePort(ctx, p)
}
