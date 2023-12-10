package port

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metrics"
)

type MetricsAdapter interface {
	ConfigureContainer(uuid uuid.UUID) error
	GetMetrics() ([]metrics.Metric, error)
}
