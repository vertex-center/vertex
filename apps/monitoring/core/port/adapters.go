package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metrics"
)

type MetricsAdapter interface {
	ConfigureContainer(uuid uuid.UUID) error
	GetMetrics(ctx context.Context) ([]metrics.Metric, error)
}
