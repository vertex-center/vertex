package reverseproxy

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/reverseproxy/adapter"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/service"
	"github.com/vertex-center/vertex/apps/reverseproxy/handler"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/log"
	"github.com/wI2L/fizz"
)

var Meta = appmeta.Meta{
	ID:          "reverse-proxy",
	Name:        "Vertex Reverse Proxy",
	Description: "Redirect traffic to your containers.",
	Icon:        "router",
	DefaultPort: "7508",
	Dependencies: []*appmeta.Meta{
		&authmeta.Meta,
	},
}

type App struct {
	ctx   *app.Context
	proxy *ProxyRouter
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

var (
	proxyService port.ProxyService
)

func (a *App) Initialize() error {
	var (
		proxyFSAdapter = adapter.NewProxyFSAdapter(nil)
	)

	proxyService = service.NewProxyService(proxyFSAdapter)
	a.proxy = NewProxyRouter(proxyService)

	go func() {
		err := a.proxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	var (
		proxyHandler = handler.NewProxyHandler(proxyService)

		redirects = r.Group("/redirects", "Redirects", "", middleware.Authenticated)
		redirect  = r.Group("/redirect", "Redirect", "", middleware.Authenticated)
	)

	redirects.GET("", []fizz.OperationOption{
		fizz.ID("getRedirects"),
		fizz.Summary("Get redirects"),
	}, proxyHandler.GetRedirects())

	redirect.POST("", []fizz.OperationOption{
		fizz.ID("addRedirect"),
		fizz.Summary("Add redirect"),
	}, proxyHandler.AddRedirect())

	redirect.DELETE("/:id", []fizz.OperationOption{
		fizz.ID("removeRedirect"),
		fizz.Summary("Remove redirect"),
	}, proxyHandler.RemoveRedirect())

	return nil
}

func (a *App) Uninitialize() error {
	return a.proxy.Stop()
}
