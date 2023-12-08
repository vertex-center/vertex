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
	"github.com/wI2L/fizz"
)

var (
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

func (a *App) Initialize() error {
	prometheusAdapter := adapter.NewMetricsPrometheusAdapter()
	metricsService = service.NewMetricsService(a.ctx, prometheusAdapter)

	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	metricsHandler := handler.NewMetricsHandler(metricsService)

	r.GET("/metrics", []fizz.OperationOption{
		fizz.ID("getMetrics"),
		fizz.Summary("Get metrics"),
	}, middleware.Authenticated, metricsHandler.Get())

	r.POST("/collector/:collector/install", []fizz.OperationOption{
		fizz.ID("installCollector"),
		fizz.Summary("Install a collector"),
	}, metricsHandler.InstallCollector())

	r.POST("/visualizer/:visualizer/install", []fizz.OperationOption{
		fizz.ID("installVisualizer"),
		fizz.Summary("Install a visualizer"),
	}, metricsHandler.InstallVisualizer())

	return nil
}
