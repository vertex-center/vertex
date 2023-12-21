package reverseproxy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

type ProxyRouter struct {
	*router.Router

	proxyService port.ProxyService
}

func NewProxyRouter(proxyService port.ProxyService) *ProxyRouter {
	r := &ProxyRouter{
		Router: router.New(nil,
			router.WithMiddleware(cors.Default()),
			router.WithMiddleware(gin.Recovery()),
		),
		proxyService: proxyService,
	}
	r.initAPIRoutes()
	return r
}

func (r *ProxyRouter) Start() error {
	proxyURL := config.Current.Addr("proxy")
	log.Info("Vertex-Proxy started", vlog.String("url", proxyURL.String()))
	return r.Router.Start(":" + proxyURL.Port())
}

func (r *ProxyRouter) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return r.Router.Stop(ctx)
}

func (r *ProxyRouter) initAPIRoutes() {
	r.Engine().Any("/*path", r.HandleProxy)
}

func (r *ProxyRouter) HandleProxy(ctx *gin.Context) {
	host := ctx.Request.Host

	redirect := r.proxyService.GetRedirectByHost(host)
	if redirect == nil {
		log.Warn("this host is not registered in the reverse proxy",
			vlog.String("host", host),
		)
		return
	}

	target, err := url.Parse(redirect.Target)
	if err != nil {
		log.Error(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, request *http.Request, err error) {
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Error(err)
		}
	}
	proxy.Director = func(request *http.Request) {
		request.Header = ctx.Request.Header
		request.Host = target.Host
		request.URL.Scheme = target.Scheme
		request.URL.Host = target.Host
		request.URL.Path = ctx.Param("path")
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
