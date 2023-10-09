package api

import (
	"context"

	"github.com/carlmjohnson/requests"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/types"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/types"
)

func GetMetrics(ctx context.Context) ([]metricstypes.Metric, *types.AppApiError) {
	var metrics []metricstypes.Metric
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Path("/api/metrics").
		ToJSON(&metrics).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return metrics, types.HandleError(err, apiError)
}

func InstallMetricsCollector(ctx context.Context, collector string) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/metrics/collector/%s/install", collector).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}

func InstallMetricsVisualizer(ctx context.Context, visualizer string) *types.AppApiError {
	var apiError types.AppApiError
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/metrics/visualizer/%s/install", visualizer).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return types.HandleError(err, apiError)
}
