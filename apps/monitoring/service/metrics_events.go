package service

import (
	"math"

	"github.com/google/uuid"
	containerstypes "github.com/vertex-center/vertex/apps/containers/types"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/types"
)

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *MetricsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case containerstypes.EventContainerStatusChange:
		s.updateStatus(e.ContainerUUID, e.ServiceID, e.Status)
	case containerstypes.EventContainerCreated:
		s.adapter.Inc(MetricIDContainersCount)
	case containerstypes.EventContainerDeleted:
		s.adapter.Dec(MetricIDContainersCount)
		s.adapter.Set(MetricIDContainerStatus, math.NaN(), e.ContainerUUID.String(), e.ServiceID)
	case containerstypes.EventContainersLoaded:
		s.adapter.Set(MetricIDContainersCount, float64(e.Count))
	}
}

func (s *MetricsService) updateStatus(uuid uuid.UUID, serviceId string, status string) {
	switch status {
	case containerstypes.ContainerStatusRunning:
		s.adapter.Set(MetricIDContainerStatus, metricstypes.MetricStatusOn, uuid.String(), serviceId)
	default:
		s.adapter.Set(MetricIDContainerStatus, metricstypes.MetricStatusOff, uuid.String(), serviceId)
	}
}
