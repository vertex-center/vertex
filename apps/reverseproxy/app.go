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
	"github.com/vertex-center/vertex/pkg/router"
)

// docapi:proxy title Vertex Reverse Proxy
// docapi:proxy description A reverse proxy manager.
// docapi:proxy version 0.0.0
// docapi:proxy filename proxy

// docapi:proxy url http://{ip}:{port-kernel}/api
// docapi:proxy urlvar ip localhost The IP address of the server.
// docapi:proxy urlvar port-kernel 7508 The port of the server.

var (
	proxyFSAdapter port.ProxyAdapter

	proxyService port.ProxyService
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
	r.Use(middleware.ReadAuth)

	proxyFSAdapter = adapter.NewProxyFSAdapter(nil)

	proxyService = service.NewProxyService(proxyFSAdapter)

	a.proxy = NewProxyRouter(proxyService)

	go func() {
		err := a.proxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	proxyHandler := handler.NewProxyHandler(proxyService)
	// docapi:proxy route /redirects vx_reverse_proxy_get_redirects
	r.GET("/redirects", middleware.Authenticated, proxyHandler.GetRedirects)
	// docapi:proxy route /redirect vx_reverse_proxy_add_redirect
	r.POST("/redirect", middleware.Authenticated, proxyHandler.AddRedirect)
	// docapi:proxy route /redirect/{id} vx_reverse_proxy_remove_redirect
	r.DELETE("/redirect/:id", middleware.Authenticated, proxyHandler.RemoveRedirect)

	return nil
}

func (a *App) Uninitialize() error {
	return a.proxy.Stop()
}
