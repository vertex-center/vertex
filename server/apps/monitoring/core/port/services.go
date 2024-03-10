package port

import (
	"context"

	containerstypes "github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/apps/monitoring/core/types"
)

type (
	MetricsService interface {
		GetCollector(ctx context.Context, collector string) (types.Collector, error)
		ConfigureVisualizer(inst *containerstypes.Container) error
		InstallVisualizer(ctx context.Context, visualizer string) error
		ConfigureCollector(inst *containerstypes.Container) error
		InstallCollector(ctx context.Context, collector string) error
	}
)
