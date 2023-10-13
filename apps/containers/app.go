package containers

import (
	"github.com/vertex-center/vertex/apps/containers/router"
	apptypes "github.com/vertex-center/vertex/types/app"
)

const (
	AppRoute = "/vx-containers"
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
		ID:          "vx-containers",
		Name:        "Vertex Containers",
		Description: "Create and manage containers.",
		Icon:        "deployed_code",
	})
	app.RegisterRouter(AppRoute, a.router)

	return nil
}
