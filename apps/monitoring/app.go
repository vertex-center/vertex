package monitoring

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/monitoring/adapter"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/service"
	"github.com/vertex-center/vertex/apps/monitoring/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

// docapi:monitoring title Vertex Monitoring
// docapi:monitoring description A monitoring service for Vertex.
// docapi:monitoring version 0.0.0
// docapi:monitoring filename monitoring

// docapi:monitoring url http://{ip}:{port-kernel}/api
// docapi:monitoring urlvar ip localhost The IP address of the server.
// docapi:monitoring urlvar port-kernel 7506 The port of the server.

var (
	prometheusAdapter port.MetricsAdapter

	metricsService port.MetricsService
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

	prometheusAdapter = adapter.NewMetricsPrometheusAdapter()

	metricsService = service.NewMetricsService(a.ctx, prometheusAdapter)

	metricsHandler := handler.NewMetricsHandler(metricsService)
	// docapi:monitoring route /metrics vx_monitoring_get_metrics
	r.GET("/metrics", middleware.Authenticated, metricsHandler.Get)
	// docapi:monitoring route /collector/{collector}/install vx_monitoring_install_collector
	r.POST("/collector/:collector/install", metricsHandler.InstallCollector)
	// docapi:monitoring route /visualizer/{visualizer}/install vx_monitoring_install_visualizer
	r.POST("/visualizer/:visualizer/install", metricsHandler.InstallVisualizer)

	return nil
}
