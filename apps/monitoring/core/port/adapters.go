package port

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
)

type MetricsAdapter interface {
	ConfigureContainer(uuid uuid.UUID) error
	RegisterMetrics(metrics []types.Metric)

	Set(metricID string, value interface{}, labels ...string)
	Inc(metricID string, labels ...string)
	Dec(metricID string, labels ...string)
}
