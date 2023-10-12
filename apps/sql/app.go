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
	*app.App
	router *router.AppRouter
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *app.App) error {
	a.App = app

	c := app.Context()
	a.router = router.NewAppRouter(c)

	app.Register(AppID, AppName)
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
