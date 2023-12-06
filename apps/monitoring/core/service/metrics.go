package service

import (
	"context"

	"github.com/google/uuid"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
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

func (s *MetricsService) InstallCollector(ctx context.Context, token string, collector string) error {
	c := containersapi.NewContainersClient(token)

	serv, apiError := c.GetService(ctx, collector)
	if apiError != nil {
		return apiError.RouterError()
	}

	inst, apiError := c.InstallService(ctx, serv.ID)
	if apiError != nil {
		return apiError.RouterError()
	}

	err := s.ConfigureCollector(inst)
	if err != nil {
		return err
	}

	apiError = c.PatchContainer(ctx, inst.UUID, containerstypes.ContainerSettings{
		Tags: []string{"Vertex Monitoring", "Vertex Monitoring - Prometheus Collector"},
	})
	if apiError != nil {
		return apiError.RouterError()
	}

	return nil
}

// ConfigureCollector will configure a container to monitor the metrics of Vertex.
func (s *MetricsService) ConfigureCollector(inst *containerstypes.Container) error {
	// TODO: Enable again, but permissions are not set correctly
	// return s.adapter.ConfigureContainer(inst.UUID)
	return nil
}

func (s *MetricsService) InstallVisualizer(ctx context.Context, token string, visualizer string) error {
	c := containersapi.NewContainersClient(token)

	serv, apiError := c.GetService(ctx, visualizer)
	if apiError != nil {
		return apiError.RouterError()
	}

	inst, apiError := c.InstallService(ctx, serv.ID)
	if apiError != nil {
		return apiError.RouterError()
	}

	err := s.ConfigureVisualizer(inst)
	if err != nil {
		return err
	}

	apiError = c.PatchContainer(ctx, inst.UUID, containerstypes.ContainerSettings{
		Tags: []string{"Vertex Monitoring", "Vertex Monitoring - Grafana Visualizer"},
	})
	if apiError != nil {
		return apiError.RouterError()
	}

	return nil
}

func (s *MetricsService) ConfigureVisualizer(inst *containerstypes.Container) error {
	// TODO: Implement
	return nil
}

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *MetricsService) OnEvent(e event.Event) error {
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
	return nil
}
