package sql

import (
	"github.com/vertex-center/vertex/apps/sql/router"
	apptypes "github.com/vertex-center/vertex/core/types/app"
)

const (
	AppRoute = "/vx-sql"
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

	c := app.Context()
	a.router = router.NewAppRouter(c)

	app.Register(apptypes.Meta{
		ID:          "vx-sql",
		Name:        "Vertex SQL",
		Description: "Create and manage SQL databases.",
		Icon:        "database",
	})
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
