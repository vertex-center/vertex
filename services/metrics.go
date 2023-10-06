package services

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

type MetricsService struct {
	uuid uuid.UUID

	adapter types.MetricsAdapterPort
	events  types.EventAdapterPort
}

func NewMetricsService(adapter types.MetricsAdapterPort, eventsAdapter types.EventAdapterPort) MetricsService {
	s := MetricsService{
		uuid: uuid.New(),

		adapter: adapter,
		events:  eventsAdapter,
	}

	s.events.AddListener(&s)

	return s
}

// ConfigureInstance will configure an instance to monitor the metrics of Vertex.
func (s *MetricsService) ConfigureInstance(inst *types.Instance) error {
	return s.adapter.ConfigureInstance(inst.UUID)
}

func (s *MetricsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceStatusChange:
		s.updateStatus(e.InstanceUUID, e.Status)
	}
}

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *MetricsService) updateStatus(uuid uuid.UUID, status string) {
	switch status {
	case types.InstanceStatusRunning:
		s.adapter.UpdateInstanceStatus(uuid, types.MetricStatusOn)
	default:
		s.adapter.UpdateInstanceStatus(uuid, types.MetricStatusOff)
	}
}
