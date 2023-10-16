package port

import (
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
)

type (
	MetricsService interface {
		GetMetrics() []types.Metric
		ConfigureVisualizer(inst *containerstypes.Container) error
		ConfigureCollector(inst *containerstypes.Container) error
	}
)
