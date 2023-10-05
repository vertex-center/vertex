package router

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

type ProxyRouter struct {
	*router.Router
}

func NewProxyRouter() ProxyRouter {
	gin.SetMode(gin.ReleaseMode)

	r := ProxyRouter{
		Router: router.New(),
	}

	r.Use(cors.Default())
	r.Use(ginutils.Logger("PROXY"))
	r.Use(gin.Recovery())

	r.initAPIRoutes()

	return r
}

func (r *ProxyRouter) Start() error {
	log.Info("Vertex-Proxy started", vlog.String("url", config.Current.HostProxy))
	addr := fmt.Sprintf(":%s", config.Current.PortProxy)
	return r.Router.Start(addr)
}

func (r *ProxyRouter) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return r.Router.Stop(ctx)
}

func (r *ProxyRouter) initAPIRoutes() {
	r.Any("/*path", proxyService.HandleProxy)
}
