package instances

import (
	"github.com/vertex-center/vertex/apps/instances/router"
	"github.com/vertex-center/vertex/types"
)

const (
	AppID   = "instances"
	AppName = "Vertex Instances"
	Route   = ""
)

type App struct {
	router *router.AppRouter
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize(registry *types.AppsRegistry) error {
	app.router = router.NewAppRouter(registry.GetContext())

	registry.RegisterApp(AppID, app)
	registry.RegisterRouter(Route, app.router)

	return nil
}

func (app *App) Uninitialize(registry *types.AppsRegistry) error {
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
