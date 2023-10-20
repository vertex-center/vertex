package port

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
)

type MetricsAdapter interface {
	// ConfigureContainer configures a container to monitor the metrics of Vertex.
	ConfigureContainer(uuid uuid.UUID) error

	// RegisterMetrics registers the metrics that can be monitored.
	RegisterMetrics(metrics []types.Metric)

	Set(metricID string, value interface{}, labels ...string)
	Inc(metricID string, labels ...string)
	Dec(metricID string, labels ...string)
}
