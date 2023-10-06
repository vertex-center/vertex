package services

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

const (
	MetricIDInstanceStatus = "vertex_instance_status"
	MetricIDInstancesCount = "vertex_instances_count"
)

type MetricsService struct {
	uuid uuid.UUID

	adapter types.MetricsAdapterPort
	events  types.EventAdapterPort

	metrics []types.Metric
}

func NewMetricsService(adapter types.MetricsAdapterPort, eventsAdapter types.EventAdapterPort) MetricsService {
	metrics := []types.Metric{
		{
			ID:          MetricIDInstanceStatus,
			Name:        "Instance Status",
			Description: "The status of the instance",
			Type:        types.MetricTypeOnOff,
			Labels:      []string{"instance_uuid"},
		},
		{
			ID:          MetricIDInstancesCount,
			Name:        "Instances Count",
			Description: "The number of instances installed",
			Type:        types.MetricTypeInteger,
		},
	}

	s := MetricsService{
		uuid: uuid.New(),

		adapter: adapter,
		events:  eventsAdapter,

		metrics: metrics,
	}

	s.adapter.RegisterMetrics(metrics)

	s.events.AddListener(&s)

	return s
}

// ConfigureInstance will configure an instance to monitor the metrics of Vertex.
func (s *MetricsService) ConfigureInstance(inst *types.Instance) error {
	return s.adapter.ConfigureInstance(inst.UUID)
}

func (s *MetricsService) GetMetrics() []types.Metric {
	return s.metrics
}

func (s *MetricsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceStatusChange:
		s.updateStatus(e.InstanceUUID, e.Status)
	case types.EventInstanceCreated:
		s.adapter.Inc(MetricIDInstancesCount)
	case types.EventInstanceDeleted:
		s.adapter.Dec(MetricIDInstancesCount)
	case types.EventInstancesLoaded:
		s.adapter.Set(MetricIDInstancesCount, float64(e.Count))
	}
}

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *MetricsService) updateStatus(uuid uuid.UUID, status string) {
	switch status {
	case types.InstanceStatusRunning:
		s.adapter.Set(MetricIDInstanceStatus, types.MetricStatusOn, uuid.String())
	default:
		s.adapter.Set(MetricIDInstanceStatus, types.MetricStatusOff, uuid.String())
	}
}
