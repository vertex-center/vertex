package reverseproxy

import (
	"github.com/vertex-center/vertex/apps/reverseproxy/router"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types/app"
)

const (
	AppID    = "vx-reverse-proxy"
	AppName  = "Vertex Reverse Proxy"
	AppRoute = "/vx-reverse-proxy"
)

type App struct {
	*app.App
	router *router.AppRouter
	proxy  *router.ProxyRouter
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *app.App) error {
	a.App = app
	a.router = router.NewAppRouter()
	a.proxy = router.NewProxyRouter(a.router.GetProxyService())

	go func() {
		err := a.proxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	app.Register(AppID, AppName)
	app.RegisterRouter(AppRoute, a.router)

	return nil
}

func (a *App) Uninitialize() error {
	return a.proxy.Stop()
}
