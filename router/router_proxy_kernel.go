package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vlog"
)

var (
	proxyKernelService services.ProxyKernelService
)

type ProxyKernelRouter struct {
	server *http.Server
	engine *gin.Engine
}

func NewProxyKernelRouter() ProxyKernelRouter {
	gin.SetMode(gin.ReleaseMode)

	r := ProxyKernelRouter{
		engine: gin.New(),
	}

	r.engine.Use(cors.Default())
	r.engine.Use(ginutils.Logger("PROXY"))
	r.engine.Use(gin.Recovery())

	r.initServices()
	r.initAPIRoutes()

	return r
}

func (r *ProxyKernelRouter) Start() error {
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

func (r *ProxyKernelRouter) Stop() error {
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

func (r *ProxyKernelRouter) initServices() {
	proxyKernelService = services.NewProxyKernelService()
}

func (r *ProxyKernelRouter) initAPIRoutes() {
	r.engine.Any("/*path", proxyKernelService.HandleProxy)
}
