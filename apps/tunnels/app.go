package tunnels

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/tunnels/handler"
	"github.com/vertex-center/vertex/common/app"
	"github.com/wI2L/fizz"
)

var Meta = app.Meta{
	ID:          "tunnels",
	Name:        "Vertex Tunnels",
	Description: "Create and manage tunnels.",
	Icon:        "subway",
	DefaultPort: "7514",
	Dependencies: []*app.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
	},
}

type App struct {
	ctx *app.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *app.Context) {
	a.ctx = ctx
}

func (a *App) Meta() app.Meta {
	return Meta
}

func (a *App) Initialize() error {
	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	providerHandler := handler.NewProviderHandler()

	r.POST("/provider/:provider/install", []fizz.OperationOption{
		fizz.ID("installProvider"),
		fizz.Summary("Install a tunnel provider"),
	}, middleware.Authenticated, providerHandler.Install())

	return nil
}
