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
	updateService        services.UpdateDependenciesService
	settingsService      services.SettingsService
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
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.Path, "client", "dist"), true)))
	r.GET("/ping", handlePing)

	runnerDockerAdapter = adapter.NewRunnerDockerAdapter()
	runnerFSAdapter = adapter.NewRunnerFSAdapter()
	instanceFSAdapter = adapter.NewInstanceFSAdapter()
	instanceLogsFSAdapter = adapter.NewInstanceLogsFSAdapter()
	eventInMemoryAdapter = adapter.NewEventInMemoryAdapter()
	packageFSAdapter = adapter.NewPackageFSAdapter(nil)
	serviceFSAdapter = adapter.NewServiceFSAdapter(nil)
	proxyFSAdapter = adapter.NewProxyFSAdapter(nil)
	settingsFSAdapter = adapter.NewSettingsFSAdapter(nil)

	proxyService = services.NewProxyService(proxyFSAdapter)
	notificationsService = services.NewNotificationsService(settingsFSAdapter, eventInMemoryAdapter, instanceFSAdapter)
	instanceService = services.NewInstanceService(serviceFSAdapter, instanceFSAdapter, runnerDockerAdapter, runnerFSAdapter, instanceLogsFSAdapter, eventInMemoryAdapter)
	packageService = services.NewPackageService(packageFSAdapter)
	serviceService = services.NewServiceService(serviceFSAdapter)
	updateService = services.NewUpdateDependenciesService(about.Version)
	settingsService = services.NewSettingsService(settingsFSAdapter)

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
	addSettingsRoutes(api.Group("/settings"))

	router.engine = r

	return router
}

func (r *Router) Start(addr string) {
	go func() {
		err := proxyService.Start()
		if err != nil {
			log.Default.Error(err)
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
		log.Default.Error(err)
	}

	err = r.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Default.Info("Vertex closed")
	} else if err != nil {
		log.Default.Error(err)
	}
}

func (r *Router) Stop() {
	// TODO: Stop() must also stop handleSignals()

	instanceService.StopAll()

	err := r.server.Shutdown(context.Background())
	if err != nil {
		log.Default.Error(err)
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
		log.Default.Info("shutdown signal sent")
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
