package service

import (
	"github.com/vertex-center/vertex/apps/containers/core/types"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/core/types"
	"math"

	"github.com/google/uuid"
)

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *MetricsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventContainerStatusChange:
		s.updateStatus(e.ContainerUUID, e.ServiceID, e.Status)
	case types.EventContainerCreated:
		s.adapter.Inc(MetricIDContainersCount)
	case types.EventContainerDeleted:
		s.adapter.Dec(MetricIDContainersCount)
		s.adapter.Set(MetricIDContainerStatus, math.NaN(), e.ContainerUUID.String(), e.ServiceID)
	case types.EventContainersLoaded:
		s.adapter.Set(MetricIDContainersCount, float64(e.Count))
	}
}

func (s *MetricsService) updateStatus(uuid uuid.UUID, serviceId string, status string) {
	switch status {
	case types.ContainerStatusRunning:
		s.adapter.Set(MetricIDContainerStatus, metricstypes.MetricStatusOn, uuid.String(), serviceId)
	default:
		s.adapter.Set(MetricIDContainerStatus, metricstypes.MetricStatusOff, uuid.String(), serviceId)
	}
}
