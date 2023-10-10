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
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/apps/instances"
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/app"
	"github.com/vertex-center/vlog"
)

var (
	settingsFSAdapter   types.SettingsAdapterPort
	sshKernelApiAdapter types.SshAdapterPort

	notificationsService services.NotificationsService

	dependenciesService services.DependenciesService
	settingsService     services.SettingsService
	hardwareService     services.HardwareService
	sshService          services.SshService
)

type Router struct {
	*router.Router

	ctx          *types.VertexContext
	appsRegistry *app.AppsRegistry
}

func NewRouter(about types.About) Router {
	gin.SetMode(gin.ReleaseMode)

	ctx := types.NewVertexContext()

	r := Router{
		Router:       router.New(),
		ctx:          ctx,
		appsRegistry: app.NewAppsRegistry(ctx),
	}

	r.Use(cors.Default())
	r.Use(ginutils.ErrorHandler())
	r.Use(ginutils.Logger("MAIN"))
	r.Use(gin.Recovery())
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.Path, "client", "dist"), true)))
	r.GET("/ping", handlePing)

	r.initApps()
	r.initAdapters()
	r.initServices(about)
	r.initAPIRoutes(about)

	return r
}

func (r *Router) Start(addr string) {
	r.ctx.SendEvent(types.EventServerStart{})

	r.handleSignals()

	err := notificationsService.StartWebhook()
	if err != nil {
		log.Error(err)
	}

	url := config.Current.VertexURL()
	log.Info("Vertex started", vlog.String("url", url))

	err = r.Router.Start(addr)
	if errors.Is(err, http.ErrServerClosed) {
		log.Info("Vertex closed")
	} else if err != nil {
		log.Error(err)
	}
}

func (r *Router) Stop() {
	// TODO: Stop() must also stop handleSignals()

	r.ctx.SendEvent(types.EventServerStop{})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.Router.Stop(ctx)
	if err != nil {
		log.Error(err)
		return
	}
}

func handlePing(c *router.Context) {
	c.JSON(gin.H{
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
		os.Exit(0)
	}()
}

func (r *Router) initApps() {
	apps := []app.App{
		sql.NewApp(),
		tunnels.NewApp(),
		monitoring.NewApp(),
		instances.NewApp(),
		reverseproxy.NewApp(),
	}
	for _, app := range apps {
		log.Info("initializing app", vlog.String("name", app.Name()))
		err := app.Initialize(r.appsRegistry)
		if err != nil {
			log.Error(err)
		}
	}
}

func (r *Router) initAdapters() {
	settingsFSAdapter = adapter.NewSettingsFSAdapter(nil)
	sshKernelApiAdapter = adapter.NewSshKernelApiAdapter()
}

func (r *Router) initServices(about types.About) {
	notificationsService = services.NewNotificationsService(r.ctx, settingsFSAdapter)
	dependenciesService = services.NewDependenciesService(r.ctx, about.Version)
	settingsService = services.NewSettingsService(settingsFSAdapter)
	hardwareService = services.NewHardwareService()
	sshService = services.NewSshService(sshKernelApiAdapter)
}

func (r *Router) initAPIRoutes(about types.About) {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, router.Error{
			Code:          "resource_not_found",
			PublicMessage: "Resource not found.",
		})
	})

	api := r.Group("/api")
	api.GET("/ping", handlePing)
	api.GET("/about", func(c *router.Context) {
		c.JSON(about)
	})

	addDependenciesRoutes(api.Group("/dependencies"))
	addSettingsRoutes(api.Group("/settings"))
	addHardwareRoutes(api.Group("/hardware"))
	addSecurityRoutes(api.Group("/security"))

	for group, r := range r.appsRegistry.GetRouters() {
		r.AddRoutes(api.Group("/app" + group))
	}
}
