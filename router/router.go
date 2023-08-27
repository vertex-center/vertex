package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex-core-golang/router/middleware"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/repository"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types"
)

var (
	runnerDockerRepo  repository.RunnerDockerRepository
	runnerFSRepo      repository.RunnerFSRepository
	instanceRepo      repository.InstanceFSRepository
	instanceLogsRepo  repository.InstanceLogsFSRepository
	eventInMemoryRepo repository.EventInMemoryRepository
	packageRepo       repository.PackageFSRepository
	serviceRepo       repository.ServiceFSRepository
	proxyRepo         repository.ProxyFSRepository

	packageService  services.PackageService
	serviceService  services.ServiceService
	proxyService    services.ProxyService
	instanceService services.InstanceService
	updateService   services.UpdateDependenciesService
)

type Router struct {
	server *http.Server
	engine *gin.Engine
}

func NewRouter(about types.About) Router {
	gin.SetMode(gin.ReleaseMode)

	router := Router{}

	r := gin.New()
	r.Use(cors.Default())
	r.Use(ginutils.Logger("MAIN"))
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorMiddleware())
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.PathClient, "dist"), true)))
	r.GET("/ping", handlePing)

	runnerDockerRepo = repository.NewRunnerDockerRepository()
	runnerFSRepo = repository.NewRunnerFSRepository()
	instanceRepo = repository.NewInstanceFSRepository()
	instanceLogsRepo = repository.NewInstanceLogsFSRepository()
	eventInMemoryRepo = repository.NewEventInMemoryRepository()
	packageRepo = repository.NewPackageFSRepository(nil)
	serviceRepo = repository.NewServiceFSRepository(nil)
	proxyRepo = repository.NewProxyFSRepository(nil)

	proxyService = services.NewProxyService(&proxyRepo)
	instanceService = services.NewInstanceService(&serviceRepo, &instanceRepo, &runnerDockerRepo, &runnerFSRepo, &instanceLogsRepo, &eventInMemoryRepo)
	packageService = services.NewPackageService(&packageRepo)
	serviceService = services.NewServiceService(&serviceRepo)
	updateService = services.NewUpdateDependenciesService(about.Version)

	api := r.Group("/api")
	api.GET("/ping", handlePing)
	api.GET("/about", func(c *gin.Context) {
		c.JSON(http.StatusOK, about)
	})

	addServicesRoutes(api.Group("/services"))
	addInstancesRoutes(api.Group("/instances"))
	addInstanceRoutes(api.Group("/instance/:instance_uuid"))
	addPackagesRoutes(api.Group("/packages"))
	addProxyRoutes(api.Group("/proxy"))
	addUpdatesRoutes(api.Group("/updates"))

	router.engine = r

	return router
}

func (r *Router) Start(addr string) {
	go func() {
		err := proxyService.Start()
		if err != nil {
			logger.Error(err).Print()
			return
		}

		instanceService.StartAll()
	}()

	r.handleSignals()

	r.server = &http.Server{
		Addr:    addr,
		Handler: r.engine,
	}

	err := r.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.Log("Vertex closed").Print()
	} else if err != nil {
		logger.Error(err).Print()
	}
}

func (r *Router) Stop() {
	// TODO: Stop() must also stop handleSignals()

	instanceService.StopAll()

	err := r.server.Shutdown(context.Background())
	if err != nil {
		logger.Error(err).Print()
		return
	}

	r.server = nil
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (r *Router) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Log("shutdown signal sent").Print()
		r.Stop()
		os.Exit(0)
	}()
}

func headersSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
