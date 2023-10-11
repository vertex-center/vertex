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
	*app.App
	router *router.AppRouter
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *app.App) error {
	a.App = app
	a.router = router.NewAppRouter(app.Context())

	app.Register(AppID, AppName)
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
