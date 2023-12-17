package port

import (
	"context"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
)

type MetricsAdapter interface {
	ConfigureContainer(uuid uuid.UUID) error
	GetMetrics(ctx context.Context) ([]metric.Metric, error)
}
