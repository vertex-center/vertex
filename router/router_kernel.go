package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

var (
	dockerCliAdapter types.DockerAdapterPort

	dockerKernelService services.DockerKernelService
	sshKernelService    services.SSHKernelService
)

type KernelRouter struct {
	server *http.Server
	engine *gin.Engine
}

func NewKernelRouter() KernelRouter {
	gin.SetMode(gin.ReleaseMode)

	r := KernelRouter{
		engine: gin.New(),
	}

	r.engine.Use(ginutils.ErrorHandler())
	r.engine.Use(ginutils.Logger("KERNEL"))
	r.engine.Use(gin.Recovery())
	r.engine.GET("/ping", handlePing)

	r.initAdapters()
	r.initServices()
	r.initAPIRoutes()

	return r
}

func (r *KernelRouter) Start() error {
	r.server = &http.Server{
		Addr:    ":6131",
		Handler: r.engine,
	}

	url := "http://localhost:6131"
	log.Info("Vertex-Kernel started",
		vlog.String("url", url),
	)

	return r.server.ListenAndServe()
}

func (r *KernelRouter) Stop() error {
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

func (r *KernelRouter) initAdapters() {
	dockerCliAdapter = adapter.NewDockerCliAdapter()
}

func (r *KernelRouter) initServices() {
	dockerKernelService = services.NewDockerKernelService(dockerCliAdapter)
	sshKernelService = services.NewSSHKernelService(nil)
}

func (r *KernelRouter) initAPIRoutes() {
	api := r.engine.Group("/api")

	addDockerKernelRoutes(api.Group("/docker"))
	addSecurityKernelRoutes(api.Group("/security"))
}
