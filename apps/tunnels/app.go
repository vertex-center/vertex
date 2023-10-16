package tunnels

import (
	"github.com/vertex-center/vertex/apps/tunnels/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

const (
	AppRoute = "/vx-tunnels"
)

type App struct {
	*apptypes.App
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *apptypes.App) error {
	a.App = app

	app.Register(apptypes.Meta{
		ID:          "vx-tunnels",
		Name:        "Vertex Tunnels",
		Description: "Create and manage tunnels.",
		Icon:        "subway",
	})

	app.RegisterRoutes(AppRoute, func(r *router.Group) {
		providerHandler := handler.NewProviderHandler()
		r.POST("/provider/:provider/install", providerHandler.Install)
	})

	return nil
}
