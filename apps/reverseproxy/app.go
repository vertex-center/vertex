package reverseproxy

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/reverseproxy/adapter"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/service"
	"github.com/vertex-center/vertex/apps/reverseproxy/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/wI2L/fizz"
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

var (
	proxyService port.ProxyService
)

func (a *App) Initialize() error {
	var (
		proxyFSAdapter = adapter.NewProxyFSAdapter(nil)
		proxyService   = service.NewProxyService(proxyFSAdapter)
	)

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

	proxyHandler := handler.NewProxyHandler(proxyService)

	r.GET("/redirects", []fizz.OperationOption{
		fizz.ID("getRedirects"),
		fizz.Summary("Get redirects"),
	}, middleware.Authenticated, proxyHandler.GetRedirects())

	r.POST("/redirect", []fizz.OperationOption{
		fizz.ID("addRedirect"),
		fizz.Summary("Add redirect"),
	}, middleware.Authenticated, proxyHandler.AddRedirect())

	r.DELETE("/redirect/:id", []fizz.OperationOption{
		fizz.ID("removeRedirect"),
		fizz.Summary("Remove redirect"),
	}, middleware.Authenticated, proxyHandler.RemoveRedirect())

	return nil
}

func (a *App) Uninitialize() error {
	return a.proxy.Stop()
}
