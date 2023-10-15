package tunnels

import (
	"github.com/vertex-center/vertex/apps/tunnels/router"
	apptypes "github.com/vertex-center/vertex/core/types/app"
)

const (
	AppRoute = "/vx-tunnels"
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
	a.router = router.NewAppRouter()

	app.Register(apptypes.Meta{
		ID:          "vx-tunnels",
		Name:        "Vertex Tunnels",
		Description: "Create and manage tunnels.",
		Icon:        "subway",
	})
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
