package monitoring

import (
	"github.com/vertex-center/vertex/apps/monitoring/adapter"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/service"
	"github.com/vertex-center/vertex/apps/monitoring/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

const (
	AppRoute = "/vx-monitoring"
)

var (
	prometheusAdapter port.MetricsAdapter

	metricsService port.MetricsService
)

type App struct {
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return apptypes.Meta{
		ID:          "vx-monitoring",
		Name:        "Vertex Monitoring",
		Description: "Create and manage containers.",
		Icon:        "monitoring",
	}
}

func (a *App) Initialize(r *router.Group) error {
	prometheusAdapter = adapter.NewMetricsPrometheusAdapter()

	metricsService = service.NewMetricsService(a.ctx, prometheusAdapter)

	metricsHandler := handler.NewMetricsHandler(metricsService)
	r.GET("/metrics", metricsHandler.Get)
	r.POST("/collector/:collector/install", metricsHandler.InstallCollector)
	r.POST("/visualizer/:visualizer/install", metricsHandler.InstallVisualizer)

	return nil
}
