package monitoring

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/monitoring/adapter"
	"github.com/vertex-center/vertex/apps/monitoring/core/service"
	"github.com/vertex-center/vertex/apps/monitoring/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

var Meta = apptypes.Meta{
	ID:          "monitoring",
	Name:        "Vertex Monitoring",
	Description: "Create and manage containers.",
	Icon:        "monitoring",
	DefaultPort: "7506",
	Dependencies: []*apptypes.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
	},
}

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
	return Meta
}

func (a *App) Initialize(r *router.Group) error {
	r.Use(middleware.ReadAuth)

	var (
		prometheusAdapter = adapter.NewMetricsPrometheusAdapter()
		metricsService    = service.NewMetricsService(a.ctx, prometheusAdapter)
		metricsHandler    = handler.NewMetricsHandler(metricsService)
	)

	r.GET("/metrics", metricsHandler.GetInfo(), middleware.Authenticated, metricsHandler.Get)
	r.POST("/collector/:collector/install", metricsHandler.InstallCollectorInfo(), metricsHandler.InstallCollector)
	r.POST("/visualizer/:visualizer/install", metricsHandler.InstallVisualizerInfo(), metricsHandler.InstallVisualizer)

	return nil
}
