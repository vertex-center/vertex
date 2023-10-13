package reverseproxy

import (
	"github.com/vertex-center/vertex/apps/reverseproxy/router"
	"github.com/vertex-center/vertex/pkg/log"
	apptypes "github.com/vertex-center/vertex/types/app"
)

const (
	AppRoute = "/vx-reverse-proxy"
)

type App struct {
	*apptypes.App
	router *router.AppRouter
	proxy  *router.ProxyRouter
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *apptypes.App) error {
	a.App = app
	a.router = router.NewAppRouter()
	a.proxy = router.NewProxyRouter(a.router.GetProxyService())

	go func() {
		err := a.proxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	app.Register(apptypes.Meta{
		ID:          "vx-reverse-proxy",
		Name:        "Vertex Reverse Proxy",
		Description: "Redirect traffic to your containers.",
		Icon:        "router",
	})
	app.RegisterRouter(AppRoute, a.router)

	return nil
}

func (a *App) Uninitialize() error {
	return a.proxy.Stop()
}
