package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/server/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/server/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/server/pkg/router"
)

type metricsHandler struct {
	metricsService port.MetricsService
}

func NewMetricsHandler(metricsService port.MetricsService) port.MetricsHandler {
	return &metricsHandler{
		metricsService: metricsService,
	}
}

type GetCollectorParams struct {
	Collector string `path:"collector"`
}

func (r *metricsHandler) GetCollector() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *GetCollectorParams) (types.Collector, error) {
		return r.metricsService.GetCollector(ctx, params.Collector)
	})
}

type InstallCollectorParams struct {
	Collector string `path:"collector"`
}

func (r *metricsHandler) InstallCollector() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *InstallCollectorParams) error {
		return r.metricsService.InstallCollector(ctx, params.Collector)
	})
}

type InstallVisualizerParams struct {
	Visualizer string `path:"visualizer"`
}

func (r *metricsHandler) InstallVisualizer() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *InstallVisualizerParams) error {
		return r.metricsService.InstallVisualizer(ctx, params.Visualizer)
	})
}
