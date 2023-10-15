package service

import (
	"github.com/google/uuid"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/adapter"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/core/types/app"
)

const (
	MetricIDContainerStatus = "vertex_container_status"
	MetricIDContainersCount = "vertex_containers_count"
)

type MetricsService struct {
	uuid    uuid.UUID
	adapter port.MetricsAdapter
	metrics []metricstypes.Metric
}

func NewMetricsService(ctx *app.Context) *MetricsService {
	metrics := []metricstypes.Metric{
		{
			ID:          MetricIDContainerStatus,
			Name:        "Container Status",
			Description: "The status of the container",
			Type:        metricstypes.MetricTypeOnOff,
			Labels:      []string{"uuid", "service_id"},
		},
		{
			ID:          MetricIDContainersCount,
			Name:        "Containers Count",
			Description: "The number of containers installed",
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

// ConfigureCollector will configure an container to monitor the metrics of Vertex.
func (s *MetricsService) ConfigureCollector(inst *containerstypes.Container) error {
	return s.adapter.ConfigureContainer(inst.UUID)
}

func (s *MetricsService) ConfigureVisualizer(inst *containerstypes.Container) error {
	// TODO: Implement
	return nil
}

func (s *MetricsService) GetMetrics() []metricstypes.Metric {
	return s.metrics
}
