package api

import (
	"context"

	metricstypes "github.com/vertex-center/vertex/apps/monitoring/core/types"
)

func (c *Client) GetMetrics(ctx context.Context) ([]metricstypes.Metric, error) {
	var metrics []metricstypes.Metric
	err := c.Request().
		Path("./metrics").
		ToJSON(&metrics).
		Fetch(ctx)
	return metrics, err
}

func (c *Client) InstallCollector(ctx context.Context, collector string) error {
	return c.Request().
		Pathf("./collector/%s/install", collector).
		Post().
		Fetch(ctx)
}

func (c *Client) InstallVisualizer(ctx context.Context, visualizer string) error {
	return c.Request().
		Pathf("./visualizer/%s/install", visualizer).
		Post().
		Fetch(ctx)
}
