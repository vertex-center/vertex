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
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return apptypes.Meta{
		ID:          "vx-tunnels",
		Name:        "Vertex Tunnels",
		Description: "Create and manage tunnels.",
		Icon:        "subway",
	}
}

func (a *App) Initialize(r *router.Group) error {
	providerHandler := handler.NewProviderHandler()
	r.POST("/provider/:provider/install", providerHandler.Install)

	return nil
}
