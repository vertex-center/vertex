package service

import (
	"math"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
	apptypes "github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/pkg/event"
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
		s.updateStatus(e.ContainerUUID, e.ServiceID, e.Status)
	case types.EventContainerCreated:
		s.metricsRegistry.Inc(MetricIDContainersCount)
	case types.EventContainerDeleted:
		s.metricsRegistry.Dec(MetricIDContainersCount)
		s.metricsRegistry.Set(MetricIDContainerStatus, math.NaN(), e.ContainerUUID.String(), e.ServiceID)
	case types.EventContainersLoaded:
		s.metricsRegistry.Set(MetricIDContainersCount, float64(e.Count))
	}
	return nil
}

func (s *metricsService) updateStatus(uuid uuid.UUID, serviceId string, status string) {
	switch status {
	case types.ContainerStatusRunning:
		s.metricsRegistry.Set(MetricIDContainerStatus, metric.StatusOn, uuid.String(), serviceId)
	default:
		s.metricsRegistry.Set(MetricIDContainerStatus, metric.StatusOff, uuid.String(), serviceId)
	}
}

func (s *metricsService) Register() error {
	return s.metricsRegistry.Register([]metric.Metric{
		{
			ID:          MetricIDContainerStatus,
			Name:        "Container Status",
			Description: "The status of the container",
			Type:        metric.TypeGauge,
			Labels:      []string{"uuid", "service_id"},
		},
		{
			ID:          MetricIDContainersCount,
			Name:        "Containers Count",
			Description: "The number of containers installed",
			Type:        metric.TypeGauge,
		},
	})
}
