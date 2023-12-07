package reverseproxy

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/reverseproxy/adapter"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/service"
	"github.com/vertex-center/vertex/apps/reverseproxy/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

var Meta = apptypes.Meta{
	ID:          "reverse-proxy",
	Name:        "Vertex Reverse Proxy",
	Description: "Redirect traffic to your containers.",
	Icon:        "router",
	DefaultPort: "7508",
	Dependencies: []*apptypes.Meta{
		&authmeta.Meta,
	},
}

type App struct {
	ctx   *apptypes.Context
	proxy *ProxyRouter
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

func (a *App) Initialize(r *router.Group) error {
	r.Use(middleware.ReadAuth())

	var (
		proxyFSAdapter = adapter.NewProxyFSAdapter(nil)
		proxyService   = service.NewProxyService(proxyFSAdapter)
		proxyHandler   = handler.NewProxyHandler(proxyService)
	)

	a.proxy = NewProxyRouter(proxyService)

	go func() {
		err := a.proxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	r.GET("/redirects", proxyHandler.GetRedirectsInfo(), middleware.Authenticated(), proxyHandler.GetRedirects())
	r.POST("/redirect", proxyHandler.AddRedirectInfo(), middleware.Authenticated(), proxyHandler.AddRedirect())
	r.DELETE("/redirect/:id", proxyHandler.RemoveRedirectInfo(), middleware.Authenticated(), proxyHandler.RemoveRedirect())

	return nil
}

func (a *App) Uninitialize() error {
	return a.proxy.Stop()
}
