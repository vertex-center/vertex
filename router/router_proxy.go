package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

type ProxyRouter struct {
	server *http.Server
	engine *gin.Engine
}

func NewProxyRouter() ProxyRouter {
	gin.SetMode(gin.ReleaseMode)

	r := ProxyRouter{
		engine: gin.New(),
	}

	r.engine.Use(cors.Default())
	r.engine.Use(ginutils.Logger("PROXY"))
	r.engine.Use(gin.Recovery())

	r.initAPIRoutes()

	return r
}

func (r *ProxyRouter) Start() error {
	r.server = &http.Server{
		Addr:    ":80",
		Handler: r.engine,
	}

	url := "http://localhost:80"
	log.Info("Vertex-Proxy started",
		vlog.String("url", url),
	)

	return r.server.ListenAndServe()
}

func (r *ProxyRouter) Stop() error {
	if r.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	r.server = nil
	return nil
}

func (r *ProxyRouter) initAPIRoutes() {
	r.engine.Any("/*path", proxyService.HandleProxy)
}
