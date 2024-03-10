package tunnels

import (
	authmeta "github.com/vertex-center/vertex/server/apps/auth/meta"
	"github.com/vertex-center/vertex/server/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/server/apps/containers/meta"
	logsmeta "github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/apps/tunnels/handler"
	"github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/common/app/appmeta"
	"github.com/wI2L/fizz"
)

var Meta = appmeta.Meta{
	ID:          "tunnels",
	Name:        "Vertex Tunnels",
	Description: "Create and manage tunnels.",
	Icon:        "subway",
	DefaultPort: "7514",
	Dependencies: []*appmeta.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
		&logsmeta.Meta,
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

func (a *App) Meta() appmeta.Meta {
	return Meta
}

func (a *App) Initialize() error {
	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	providerHandler := handler.NewProviderHandler()
	provider := r.Group("/provider/:provider", "Provider", "", middleware.Authenticated)

	provider.POST("/install", []fizz.OperationOption{
		fizz.ID("installProvider"),
		fizz.Summary("Install a tunnel provider"),
	}, providerHandler.Install())

	return nil
}
