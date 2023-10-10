package service

import (
	"math"

	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/apps/monitoring/adapter"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/types"
)

const (
	MetricIDInstanceStatus = "vertex_instance_status"
	MetricIDInstancesCount = "vertex_instances_count"
)

type MetricsService struct {
	adapter metricstypes.MetricsAdapterPort
	metrics []metricstypes.Metric
}

func NewMetricsService() *MetricsService {
	metrics := []metricstypes.Metric{
		{
			ID:          MetricIDInstanceStatus,
			Name:        "Instance Status",
			Description: "The status of the instance",
			Type:        metricstypes.MetricTypeOnOff,
			Labels:      []string{"uuid", "service_id"},
		},
		{
			ID:          MetricIDInstancesCount,
			Name:        "Instances Count",
			Description: "The number of instances installed",
			Type:        metricstypes.MetricTypeInteger,
		},
	}

	s := &MetricsService{
		adapter: adapter.NewMetricsPrometheusAdapter(),
		metrics: metrics,
	}

	s.adapter.RegisterMetrics(metrics)

	return s
}

// ConfigureCollector will configure an instance to monitor the metrics of Vertex.
func (s *MetricsService) ConfigureCollector(inst *instancestypes.Instance) error {
	return s.adapter.ConfigureInstance(inst.UUID)
}

func (s *MetricsService) ConfigureVisualizer(inst *instancestypes.Instance) error {
	// TODO: Implement
	return nil
}

func (s *MetricsService) GetMetrics() []metricstypes.Metric {
	return s.metrics
}

func (s *MetricsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case instancestypes.EventInstanceStatusChange:
		s.updateStatus(e.InstanceUUID, e.ServiceID, e.Status)
	case instancestypes.EventInstanceCreated:
		s.adapter.Inc(MetricIDInstancesCount)
	case instancestypes.EventInstanceDeleted:
		s.adapter.Dec(MetricIDInstancesCount)
		s.adapter.Set(MetricIDInstanceStatus, math.NaN(), e.InstanceUUID.String(), e.ServiceID)
	case instancestypes.EventInstancesLoaded:
		s.adapter.Set(MetricIDInstancesCount, float64(e.Count))
	}
}

func (s *MetricsService) updateStatus(uuid uuid.UUID, serviceId string, status string) {
	switch status {
	case instancestypes.InstanceStatusRunning:
		s.adapter.Set(MetricIDInstanceStatus, metricstypes.MetricStatusOn, uuid.String(), serviceId)
	default:
		s.adapter.Set(MetricIDInstanceStatus, metricstypes.MetricStatusOff, uuid.String(), serviceId)
	}
}
