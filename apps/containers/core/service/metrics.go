package service

import (
	"math"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	monitoringtypes "github.com/vertex-center/vertex/apps/monitoring/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
)

const (
	MetricIDContainerStatus = "vertex_container_status"
	MetricIDContainersCount = "vertex_containers_count"
)

type MetricsService struct {
	uuid uuid.UUID
	ctx  *apptypes.Context
}

func NewMetricsService(ctx *apptypes.Context) port.MetricsService {
	s := &MetricsService{
		uuid: uuid.New(),
		ctx:  ctx,
	}
	ctx.AddListener(s)
	return s
}

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *MetricsService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case vtypes.EventServerStart:
		return s.ctx.DispatchEvent(monitoringtypes.EventRegisterMetrics{
			Metrics: []monitoringtypes.Metric{
				{
					ID:          MetricIDContainerStatus,
					Name:        "Container Status",
					Description: "The status of the container",
					Type:        monitoringtypes.MetricTypeOnOff,
					Labels:      []string{"uuid", "service_id"},
				},
				{
					ID:          MetricIDContainersCount,
					Name:        "Containers Count",
					Description: "The number of containers installed",
					Type:        monitoringtypes.MetricTypeInteger,
				},
			},
		})
	case types.EventContainerStatusChange:
		return s.updateStatus(e.ContainerUUID, e.ServiceID, e.Status)
	case types.EventContainerCreated:
		return s.ctx.DispatchEvent(monitoringtypes.EventIncrementMetric{
			MetricID: MetricIDContainersCount,
		})
	case types.EventContainerDeleted:
		err := s.ctx.DispatchEvent(monitoringtypes.EventDecrementMetric{
			MetricID: MetricIDContainersCount,
		})
		if err != nil {
			return err
		}
		return s.ctx.DispatchEvent(monitoringtypes.EventSetMetric{
			MetricID: MetricIDContainerStatus,
			Value:    math.NaN(),
			Labels:   []string{e.ContainerUUID.String(), e.ServiceID},
		})
	case types.EventContainersLoaded:
		return s.ctx.DispatchEvent(monitoringtypes.EventSetMetric{
			MetricID: MetricIDContainersCount,
			Value:    float64(e.Count),
		})
	}
	return nil
}

func (s *MetricsService) updateStatus(uuid uuid.UUID, serviceId string, status string) error {
	switch status {
	case types.ContainerStatusRunning:
		return s.ctx.DispatchEvent(monitoringtypes.EventSetMetric{
			MetricID: MetricIDContainerStatus,
			Value:    monitoringtypes.MetricStatusOn,
			Labels:   []string{uuid.String(), serviceId},
		})
	default:
		return s.ctx.DispatchEvent(monitoringtypes.EventSetMetric{
			MetricID: MetricIDContainerStatus,
			Value:    monitoringtypes.MetricStatusOff,
			Labels:   []string{uuid.String(), serviceId},
		})
	}
}
