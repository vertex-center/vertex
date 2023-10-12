package instances

import (
	"github.com/vertex-center/vertex/apps/instances/router"
	apptypes "github.com/vertex-center/vertex/types/app"
)

const (
	AppRoute = "/vx-instances"
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
		ID:          "vx-instances",
		Name:        "Vertex Instances",
		Description: "Create and manage instances.",
		Icon:        "storage",
	})
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
