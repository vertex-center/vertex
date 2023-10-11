package service

import (
	"math"

	"github.com/google/uuid"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/types"
)

func (s *MetricsService) GetUUID() uuid.UUID {
	return s.uuid
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
