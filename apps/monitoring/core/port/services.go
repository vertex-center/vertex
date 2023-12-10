package port

import (
	"context"

	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
)

type (
	MetricsService interface {
		GetCollector(ctx context.Context, collector string) (types.Collector, error)
		ConfigureVisualizer(inst *containerstypes.Container) error
		InstallVisualizer(ctx context.Context, token string, visualizer string) error
		ConfigureCollector(inst *containerstypes.Container) error
		InstallCollector(ctx context.Context, token string, collector string) error
	}
)
