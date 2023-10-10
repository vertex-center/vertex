package instances

import (
	"github.com/vertex-center/vertex/apps/instances/router"
	"github.com/vertex-center/vertex/types/app"
)

const (
	AppID    = "instances"
	AppName  = "Vertex Instances"
	AppRoute = "/vx-instances"
)

type App struct {
	router *router.AppRouter
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize(registry *app.AppsRegistry) error {
	app.router = router.NewAppRouter(registry.GetContext())

	registry.RegisterApp(AppID, app)
	registry.RegisterRouter(AppRoute, app.router)

	return nil
}

func (app *App) Uninitialize(registry *app.AppsRegistry) error {
	registry.UnregisterApp(AppID)
	registry.UnregisterRouter(AppRoute)

	return nil
}

func (app *App) Name() string {
	return AppName
}

func (app *App) OnEvent(e interface{}) {
	for _, s := range app.router.GetServices() {
		s.OnEvent(e)
	}
}
