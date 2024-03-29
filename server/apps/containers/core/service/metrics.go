package service

import (
	"math"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/apps/monitoring/core/types/metric"
	apptypes "github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/pkg/event"
)

const (
	MetricIDContainerStatus = "vertex_container_status"
	MetricIDContainersCount = "vertex_containers_count"
)

type metricsService struct {
	uuid            uuid.UUID
	ctx             *apptypes.Context
	metricsRegistry *metric.Registry
}

func NewMetricsService(ctx *apptypes.Context) port.MetricsService {
	s := &metricsService{
		uuid:            uuid.New(),
		ctx:             ctx,
		metricsRegistry: metric.NewServer(),
	}
	ctx.AddListener(s)
	err := s.Register()
	if err != nil {
		panic(err)
	}
	return s
}

func (s *metricsService) GetRegistry() *metric.Registry {
	return s.metricsRegistry
}

func (s *metricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *metricsService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case types.EventContainerStatusChange:
		s.updateStatus(e.ContainerID, e.Status)
	case types.EventContainerCreated:
		s.metricsRegistry.Inc(MetricIDContainersCount)
	case types.EventContainerDeleted:
		s.metricsRegistry.Dec(MetricIDContainersCount)
		s.metricsRegistry.Set(MetricIDContainerStatus, math.NaN(), e.ContainerID.String())
	case types.EventContainersLoaded:
		s.metricsRegistry.Set(MetricIDContainersCount, float64(e.Count))
	}
	return nil
}

func (s *metricsService) updateStatus(uuid uuid.UUID, status string) {
	switch status {
	case types.ContainerStatusRunning:
		s.metricsRegistry.Set(MetricIDContainerStatus, metric.StatusOn, uuid.String())
	default:
		s.metricsRegistry.Set(MetricIDContainerStatus, metric.StatusOff, uuid.String())
	}
}

func (s *metricsService) Register() error {
	return s.metricsRegistry.Register([]metric.Metric{
		{
			ID:          MetricIDContainerStatus,
			Name:        "Container Status",
			Description: "The status of the container",
			Type:        metric.TypeGauge,
			Labels:      []string{"uuid"},
		},
		{
			ID:          MetricIDContainersCount,
			Name:        "Containers Count",
			Description: "The number of containers installed",
			Type:        metric.TypeGauge,
		},
	})
}
