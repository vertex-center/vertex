package api

import (
	"context"

	metricstypes "github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/core/types/api"
)

func (c *Client) GetMetrics(ctx context.Context) ([]metricstypes.Metric, *api.Error) {
	var metrics []metricstypes.Metric
	var apiError api.Error
	err := c.Request().
		Path("./metrics").
		ToJSON(&metrics).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return metrics, api.HandleError(err, apiError)
}

func (c *Client) InstallCollector(ctx context.Context, collector string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./collector/%s/install", collector).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) InstallVisualizer(ctx context.Context, visualizer string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./visualizer/%s/install", visualizer).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}
