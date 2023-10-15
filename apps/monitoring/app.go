package monitoring

import (
	"github.com/vertex-center/vertex/apps/monitoring/router"
	apptypes "github.com/vertex-center/vertex/core/types/app"
)

const (
	AppRoute = "/vx-monitoring"
)

type App struct {
	*apptypes.App
	router *router.AppRouter
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *apptypes.App) error {
	a.App = app
	a.router = router.NewAppRouter(app.Context())

	app.Register(apptypes.Meta{
		ID:          "vx-monitoring",
		Name:        "Vertex Monitoring",
		Description: "Create and manage containers.",
		Icon:        "monitoring",
	})
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
