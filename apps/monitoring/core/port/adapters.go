package port

import (
	"context"

	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
	"github.com/vertex-center/vertex/common/uuid"
)

type MetricsAdapter interface {
	ConfigureContainer(uuid uuid.UUID) error
	GetMetrics(ctx context.Context) ([]metric.Metric, error)
}
