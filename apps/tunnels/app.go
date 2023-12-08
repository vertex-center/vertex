package tunnels

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/tunnels/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/wI2L/fizz"
)

var Meta = apptypes.Meta{
	ID:          "tunnels",
	Name:        "Vertex Tunnels",
	Description: "Create and manage tunnels.",
	Icon:        "subway",
	DefaultPort: "7514",
	Dependencies: []*apptypes.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
	},
}

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
	return Meta
}

func (a *App) Initialize(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	providerHandler := handler.NewProviderHandler()

	r.POST("/provider/:provider/install", []fizz.OperationOption{
		fizz.ID("installProvider"),
		fizz.Summary("Install a tunnel provider"),
	}, middleware.Authenticated, providerHandler.Install())

	return nil
}
