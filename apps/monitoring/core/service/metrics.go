package service

import (
	"github.com/google/uuid"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type MetricsService struct {
	uuid    uuid.UUID
	adapter port.MetricsAdapter
	metrics []types.Metric
}

func NewMetricsService(ctx *app.Context, metricsAdapter port.MetricsAdapter) port.MetricsService {
	s := &MetricsService{
		uuid:    uuid.New(),
		adapter: metricsAdapter,
		metrics: []types.Metric{},
	}
	ctx.AddListener(s)
	return s
}

func (s *MetricsService) GetMetrics() []types.Metric {
	return s.metrics
}

// ConfigureCollector will configure a container to monitor the metrics of Vertex.
func (s *MetricsService) ConfigureCollector(inst *containerstypes.Container) error {
	return s.adapter.ConfigureContainer(inst.UUID)
}

func (s *MetricsService) ConfigureVisualizer(inst *containerstypes.Container) error {
	// TODO: Implement
	return nil
}

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *MetricsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventRegisterMetrics:
		log.Info("registering metrics", vlog.Int("count", len(e.Metrics)))
		s.metrics = append(s.metrics, e.Metrics...)
		s.adapter.RegisterMetrics(e.Metrics)
	case types.EventSetMetric:
		s.adapter.Set(e.MetricID, e.Value, e.Labels...)
	case types.EventIncrementMetric:
		s.adapter.Inc(e.MetricID, e.Labels...)
	case types.EventDecrementMetric:
		s.adapter.Dec(e.MetricID, e.Labels...)
	}
}
