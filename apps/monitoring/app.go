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
	*apptypes.App
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *apptypes.App) error {
	a.App = app

	prometheusAdapter = adapter.NewMetricsPrometheusAdapter()

	metricsService = service.NewMetricsService(app.Context(), prometheusAdapter)

	app.Register(apptypes.Meta{
		ID:          "vx-monitoring",
		Name:        "Vertex Monitoring",
		Description: "Create and manage containers.",
		Icon:        "monitoring",
	})

	app.RegisterRoutes(AppRoute, func(r *router.Group) {
		metricsHandler := handler.NewMetricsHandler(metricsService)

		r.GET("/metrics", metricsHandler.Get)
		r.POST("/collector/:collector/install", metricsHandler.InstallCollector)
		r.POST("/visualizer/:visualizer/install", metricsHandler.InstallVisualizer)
	})

	return nil
}
