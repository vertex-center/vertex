package port

import (
	"context"

	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
)

type MetricsAdapter interface {
	ConfigureContainer(uuid types.ContainerID) error
	GetMetrics(ctx context.Context) ([]metric.Metric, error)
}
