package sql

import (
	"github.com/vertex-center/vertex/apps/sql/router"
	"github.com/vertex-center/vertex/types"
)

const (
	AppID   = "sql"
	AppName = "Vertex SQL"
	Route   = "/sql"
)

type App struct {
	router *router.AppRouter
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize(registry *types.AppRegistry) error {
	app.router = router.NewAppRouter()

	registry.RegisterApp(AppID, app)
	registry.RegisterRouter(Route, app.router)

	return nil
}

func (app *App) Uninitialize(registry *types.AppRegistry) error {
	registry.UnregisterApp(AppID)
	registry.UnregisterRouter(Route)

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
