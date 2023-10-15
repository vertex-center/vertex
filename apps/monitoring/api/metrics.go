package api

import (
	"context"
	"github.com/vertex-center/vertex/apps/monitoring"
	metricstypes "github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func GetMetrics(ctx context.Context) ([]metricstypes.Metric, *api.Error) {
	var metrics []metricstypes.Metric
	var apiError api.Error
	err := api.AppRequest(monitoring.AppRoute).
		Path("./metrics").
		ToJSON(&metrics).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return metrics, api.HandleError(err, apiError)
}

func InstallCollector(ctx context.Context, collector string) *api.Error {
	var apiError api.Error
	err := api.AppRequest(monitoring.AppRoute).
		Pathf("./collector/%s/install", collector).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func InstallVisualizer(ctx context.Context, visualizer string) *api.Error {
	var apiError api.Error
	err := api.AppRequest(monitoring.AppRoute).
		Pathf("./visualizer/%s/install", visualizer).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}
