package monitoring

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	logsmeta "github.com/vertex-center/vertex/apps/logs/meta"
	"github.com/vertex-center/vertex/apps/monitoring/adapter"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/service"
	"github.com/vertex-center/vertex/apps/monitoring/handler"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/wI2L/fizz"
)

var (
	metricsService port.MetricsService
)

var Meta = appmeta.Meta{
	ID:          "monitoring",
	Name:        "Vertex Monitoring",
	Description: "Create and manage containers.",
	Icon:        "monitoring",
	DefaultPort: "7506",
	Dependencies: []*appmeta.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
		&logsmeta.Meta,
	},
}

type App struct {
	ctx *app.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *app.Context) {
	a.ctx = ctx
}

func (a *App) Meta() appmeta.Meta {
	return Meta
}

func (a *App) Initialize() error {
	metricsAdapter := adapter.NewMetricsPrometheusAdapter()
	metricsService = service.NewMetricsService(metricsAdapter)
	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	var (
		metricsHandler = handler.NewMetricsHandler(metricsService)

		collector  = r.Group("/collector/:collector", "Collector", "", middleware.Authenticated)
		visualizer = r.Group("/visualizer/:visualizer", "Visualizer", "", middleware.Authenticated)
	)

	collector.GET("", []fizz.OperationOption{
		fizz.ID("getCollector"),
		fizz.Summary("Get a collector"),
	}, metricsHandler.GetCollector())

	collector.POST("/install", []fizz.OperationOption{
		fizz.ID("installCollector"),
		fizz.Summary("Install a collector"),
	}, metricsHandler.InstallCollector())

	visualizer.POST("/install", []fizz.OperationOption{
		fizz.ID("installVisualizer"),
		fizz.Summary("Install a visualizer"),
	}, metricsHandler.InstallVisualizer())

	return nil
}
