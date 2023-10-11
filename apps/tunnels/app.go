package tunnels

import (
	"github.com/vertex-center/vertex/apps/tunnels/router"
	"github.com/vertex-center/vertex/types/app"
)

const (
	AppID    = "vx-tunnels"
	AppName  = "Vertex Tunnels"
	AppRoute = "/vx-tunnels"
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
	a.router = router.NewAppRouter()

	app.Register(AppID, AppName)
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
