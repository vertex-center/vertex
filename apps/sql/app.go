package sql

import (
	"github.com/vertex-center/vertex/apps/sql/router"
	"github.com/vertex-center/vertex/types/app"
)

const (
	AppID    = "vx-sql"
	AppName  = "Vertex SQL"
	AppRoute = "/vx-sql"
)

type App struct {
	router *router.AppRouter
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize(registry *app.AppsRegistry) error {
	app.router = router.NewAppRouter()

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
