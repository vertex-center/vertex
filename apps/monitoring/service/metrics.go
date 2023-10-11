package service

import (
	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/apps/monitoring/adapter"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/types"
	"github.com/vertex-center/vertex/types/app"
)

const (
	MetricIDInstanceStatus = "vertex_instance_status"
	MetricIDInstancesCount = "vertex_instances_count"
)

type MetricsService struct {
	uuid    uuid.UUID
	adapter metricstypes.MetricsAdapterPort
	metrics []metricstypes.Metric
}

func NewMetricsService(ctx *app.Context) *MetricsService {
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
		uuid:    uuid.New(),
		adapter: adapter.NewMetricsPrometheusAdapter(),
		metrics: metrics,
	}

	s.adapter.RegisterMetrics(metrics)
	ctx.AddListener(s)

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
