package router

import (
	"context"
	"fmt"
	adapter2 "github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/core/port"
	service2 "github.com/vertex-center/vertex/core/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

var (
	dockerCliAdapter port.DockerAdapter
	sshAdapter       port.SshAdapter

	dockerKernelService service2.DockerKernelService
	sshKernelService    service2.SshKernelService
)

type KernelRouter struct {
	*router.Router
}

func NewKernelRouter() KernelRouter {
	gin.SetMode(gin.ReleaseMode)

	r := KernelRouter{
		Router: router.New(),
	}

	r.Use(ginutils.ErrorHandler())
	r.Use(ginutils.Logger("KERNEL"))
	r.Use(gin.Recovery())
	r.GET("/ping", handlePing)

	r.initAdapters()
	r.initServices()
	r.initAPIRoutes()

	return r
}

func (r *KernelRouter) Start() error {
	log.Info("vertex-kernel started", vlog.String("url", config.KernelCurrent.KernelURL()))
	addr := fmt.Sprintf(":%s", config.KernelCurrent.PortKernel)
	return r.Router.Start(addr)
}

func (r *KernelRouter) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return r.Router.Stop(ctx)
}

func (r *KernelRouter) initAdapters() {
	dockerCliAdapter = adapter2.NewDockerCliAdapter()
	sshAdapter = adapter2.NewSshFsAdapter(nil)
}

func (r *KernelRouter) initServices() {
	dockerKernelService = service2.NewDockerKernelService(dockerCliAdapter)
	sshKernelService = service2.NewSshKernelService(sshAdapter)
}

func (r *KernelRouter) initAPIRoutes() {
	api := r.Group("/api")

	addDockerKernelRoutes(api.Group("/docker"))
	addSecurityKernelRoutes(api.Group("/security"))
}
