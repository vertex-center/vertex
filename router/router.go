package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types"
)

var (
	runnerDockerAdapter   types.RunnerAdapterPort
	runnerFSAdapter       types.RunnerAdapterPort
	instanceFSAdapter     types.InstanceAdapterPort
	instanceLogsFSAdapter types.InstanceLogsAdapterPort
	eventInMemoryAdapter  types.EventAdapterPort
	packageFSAdapter      types.PackageAdapterPort
	serviceFSAdapter      types.ServiceAdapterPort
	proxyFSAdapter        types.ProxyAdapterPort
	settingsFSAdapter     types.SettingsAdapterPort

	packageService       services.PackageService
	notificationsService services.NotificationsService
	serviceService       services.ServiceService
	proxyService         services.ProxyService
	instanceService      services.InstanceService
	dependenciesService  services.DependenciesService
	settingsService      services.SettingsService
)

type Router struct {
	server *http.Server
	engine *gin.Engine
}

func NewRouter(about types.About) Router {
	gin.SetMode(gin.ReleaseMode)

	r := Router{
		engine: gin.New(),
	}

	r.engine.Use(cors.Default())
	r.engine.Use(ginutils.Logger("MAIN"))
	r.engine.Use(gin.Recovery())
	r.engine.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.Path, "client", "dist"), true)))
	r.engine.GET("/ping", handlePing)

	r.initAdapters()
	r.initServices(about)
	r.initAPIRoutes(about)

	return r
}

func (r *Router) Start(addr string) {
	go func() {
		err := proxyService.Start()
		if err != nil {
			log.Error(err)
			return
		}

		instanceService.StartAll()
	}()

	r.handleSignals()

	r.server = &http.Server{
		Addr:    addr,
		Handler: r.engine,
	}

	err := notificationsService.StartWebhook()
	if err != nil {
		log.Error(err)
	}

	err = r.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Info("Vertex closed")
	} else if err != nil {
		log.Error(err)
	}
}

func (r *Router) Stop() {
	// TODO: Stop() must also stop handleSignals()

	instanceService.StopAll()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.server.Shutdown(ctx)
	if err != nil {
		log.Error(err)
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
		log.Info("shutdown signal sent")
		r.Stop()
		err := proxyService.Stop()
		if err != nil {
			log.Error(err)
		}
		os.Exit(0)
	}()
}

func (r *Router) initAdapters() {
	runnerDockerAdapter = adapter.NewRunnerDockerAdapter()
	runnerFSAdapter = adapter.NewRunnerFSAdapter()
	instanceFSAdapter = adapter.NewInstanceFSAdapter(nil)
	instanceLogsFSAdapter = adapter.NewInstanceLogsFSAdapter()
	eventInMemoryAdapter = adapter.NewEventInMemoryAdapter()
	packageFSAdapter = adapter.NewPackageFSAdapter(nil)
	serviceFSAdapter = adapter.NewServiceFSAdapter(nil)
	proxyFSAdapter = adapter.NewProxyFSAdapter(nil)
	settingsFSAdapter = adapter.NewSettingsFSAdapter(nil)
}

func (r *Router) initServices(about types.About) {
	proxyService = services.NewProxyService(proxyFSAdapter)
	notificationsService = services.NewNotificationsService(settingsFSAdapter, eventInMemoryAdapter, instanceFSAdapter)
	instanceService = services.NewInstanceService(serviceFSAdapter, instanceFSAdapter, runnerDockerAdapter, runnerFSAdapter, instanceLogsFSAdapter, eventInMemoryAdapter)
	packageService = services.NewPackageService(packageFSAdapter)
	serviceService = services.NewServiceService(serviceFSAdapter)
	dependenciesService = services.NewDependenciesService(about.Version)
	settingsService = services.NewSettingsService(settingsFSAdapter)
}

func (r *Router) initAPIRoutes(about types.About) {
	api := r.engine.Group("/api")
	api.GET("/ping", handlePing)
	api.GET("/about", func(c *gin.Context) {
		c.JSON(http.StatusOK, about)
	})

	addServicesRoutes(api.Group("/services"))
	addInstancesRoutes(api.Group("/instances"))
	addInstanceRoutes(api.Group("/instance/:instance_uuid"))
	addPackagesRoutes(api.Group("/packages"))
	addProxyRoutes(api.Group("/proxy"))
	addDependenciesRoutes(api.Group("/dependencies"))
	addSettingsRoutes(api.Group("/settings"))
}

func headersSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
