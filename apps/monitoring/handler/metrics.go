package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type metricsHandler struct {
	metricsService port.MetricsService
}

func NewMetricsHandler(metricsService port.MetricsService) port.MetricsHandler {
	return &metricsHandler{
		metricsService: metricsService,
	}
}

func (r *metricsHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) ([]types.Metric, error) {
		return r.metricsService.GetMetrics(), nil
	})
}

func (r *metricsHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getMetrics"),
		oapi.Summary("Get metrics"),
	}
}

type InstallCollectorParams struct {
	Collector string `path:"collector"`
}

func (r *metricsHandler) InstallCollector() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InstallCollectorParams) error {
		token := c.MustGet("token").(string)
		return r.metricsService.InstallCollector(c, token, params.Collector)
	})
}

func (r *metricsHandler) InstallCollectorInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("installCollector"),
		oapi.Summary("Install a collector"),
	}
}

type InstallVisualizerParams struct {
	Visualizer string `path:"visualizer"`
}

func (r *metricsHandler) InstallVisualizer() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InstallVisualizerParams) error {
		token := c.MustGet("token").(string)
		return r.metricsService.InstallVisualizer(c, token, params.Visualizer)
	})
}

func (r *metricsHandler) InstallVisualizerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("installVisualizer"),
		oapi.Summary("Install a visualizer"),
	}
}
