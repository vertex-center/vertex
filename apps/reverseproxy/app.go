package reverseproxy

import (
	"github.com/vertex-center/vertex/apps/reverseproxy/router"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
)

const (
	AppID   = "reverse-proxy"
	AppName = "Vertex Reverse Proxy"
	Route   = "/reverse-proxy"
)

type App struct {
	router *router.AppRouter
	proxy  *router.ProxyRouter
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize(registry *types.AppsRegistry) error {
	app.router = router.NewAppRouter()
	app.proxy = router.NewProxyRouter(app.router.GetProxyService())

	go func() {
		err := app.proxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	registry.RegisterApp(AppID, app)
	registry.RegisterRouter(Route, app.router)

	return nil
}

func (app *App) Uninitialize(registry *types.AppsRegistry) error {
	err := app.proxy.Stop()

	registry.UnregisterApp(AppID)
	registry.UnregisterRouter(Route)

	if err != nil {
		return err
	}

	return nil
}

func (app *App) Name() string {
	return AppName
}

func (app *App) OnEvent(e interface{}) {
	for _, s := range app.router.GetServices() {
		s.OnEvent(e)
	}
}
