package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	adapter2 "github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/core/port"
	service2 "github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/handler"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/updates"
	"github.com/vertex-center/vlog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"time"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/migration"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
)

// version, commit and date will be overridden by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	r   *router.Router
	ctx *types.VertexContext

	settingsFSAdapter   port.SettingsAdapter
	sshKernelApiAdapter port.SshAdapter
	baselinesApiAdapter port.BaselinesAdapter

	appsService          *service2.AppsService
	notificationsService service2.NotificationsService
	settingsService      *service2.SettingsService
	hardwareService      *service2.HardwareService
	sshService           *service2.SshService
	updateService        *service2.UpdateService
)

func main() {
	defer log.Default.Close()

	log.Info("Vertex starting...")

	postMigrationCommands, err := migration.NewMigrationTool(storage.Path).Migrate()
	if err != nil {
		panic(err)
	}

	parseArgs()

	checkNotRoot()

	gin.SetMode(gin.ReleaseMode)
	ctx = types.NewVertexContext()
	r = router.New()
	r.Use(cors.Default())
	r.Use(ginutils.ErrorHandler())
	r.Use(ginutils.Logger("MAIN"))
	r.Use(gin.Recovery())

	about := types.About{
		Version: version,
		Commit:  commit,
		Date:    date,

		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	initAdapters()
	initServices(about)
	initRoutes(about)
	handleSignals()

	err = net.Wait("google.com:80")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	ctx.DispatchEvent(types.EventServerStart{
		PostMigrationCommands: postMigrationCommands,
	})

	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.Path, "client", "dist"), true)))

	err = notificationsService.StartWebhook()
	if err != nil {
		log.Error(err)
	}

	startRouter()
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Info("shutdown signal sent")
		stopRouter()
		os.Exit(0)
	}()
}

func parseArgs() {
	flagVersion := flag.Bool("version", false, "Print vertex version")
	flagV := flag.Bool("v", false, "Print vertex version")
	flagDate := flag.Bool("date", false, "Print the release date")
	flagCommit := flag.Bool("commit", false, "Print the commit hash")

	var (
		flagHost = flag.String("host", config.Current.Host, "The Vertex access url")

		flagPort           = flag.String("port", config.Current.Port, "The Vertex port")
		flagPortKernel     = flag.String("port-kernel", config.Current.PortKernel, "The Vertex Kernel port")
		flagPortProxy      = flag.String("port-proxy", config.Current.PortProxy, "The Vertex Proxy port")
		flagPortPrometheus = flag.String("port-prometheus", config.Current.PortPrometheus, "The Prometheus port")
	)

	flag.Parse()

	if *flagVersion || *flagV {
		fmt.Println(version)
		os.Exit(0)
	}
	if *flagDate {
		fmt.Println(date)
		os.Exit(0)
	}
	if *flagCommit {
		fmt.Println(commit)
		os.Exit(0)
	}
	config.Current.Host = *flagHost
	config.Current.Port = *flagPort
	config.Current.PortKernel = *flagPortKernel
	config.Current.PortProxy = *flagPortProxy
	config.Current.PortPrometheus = *flagPortPrometheus
}

func checkNotRoot() {
	if os.Getuid() == 0 {
		log.Warn("while vertex-kernel must be run as root, the vertex user should not be root")
	}
}

func initAdapters() {
	settingsFSAdapter = adapter2.NewSettingsFSAdapter(nil)
	sshKernelApiAdapter = adapter2.NewSshKernelApiAdapter()
	baselinesApiAdapter = adapter2.NewBaselinesApiAdapter()
}

func initServices(about types.About) {
	// Update service must be initialized before all other services, because it
	// is responsible for downloading dependencies for other services.
	updateService = service2.NewUpdateService(ctx, baselinesApiAdapter, []types.Updater{
		updates.NewVertexUpdater(about),
		updates.NewVertexClientUpdater(path.Join(storage.Path, "client")),
		updates.NewRepositoryUpdater("vertex_services", path.Join(storage.Path, "services"), "vertex-center", "vertex-services"),
	})
	appsService = service2.NewAppsService(ctx, r,
		[]app.Interface{
			sql.NewApp(),
			tunnels.NewApp(),
			monitoring.NewApp(),
			containers.NewApp(),
			reverseproxy.NewApp(),
		},
	)
	notificationsService = service2.NewNotificationsService(ctx, settingsFSAdapter)
	settingsService = service2.NewSettingsService(settingsFSAdapter)
	//services.NewSetupService(r.ctx)
	hardwareService = service2.NewHardwareService()
	sshService = service2.NewSshService(sshKernelApiAdapter)
}

func initRoutes(about types.About) {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, router.Error{
			Code:          "resource_not_found",
			PublicMessage: "Resource not found.",
		})
	})

	api := r.Group("/api")
	api.GET("/about", func(c *router.Context) {
		c.JSON(about)
	})

	if config.Current.Debug() {
		api.POST("/hard-reset", func(c *router.Context) {
			ctx.DispatchEvent(types.EventServerHardReset{})
			c.OK()
		})
	}

	appsHandler := handler.NewAppsHandler(appsService)
	apps := api.Group("/apps")
	apps.GET("", appsHandler.Get)

	hardwareHandler := handler.NewHardwareHandler(hardwareService)
	hardware := api.Group("/hardware")
	hardware.GET("", hardwareHandler.Get)

	updateHandler := handler.NewUpdateHandler(updateService, settingsService)
	update := api.Group("/update")
	update.GET("", updateHandler.Get)
	update.POST("", updateHandler.Install)

	settingsHandler := handler.NewSettingsHandler(settingsService)
	settings := api.Group("/settings")
	settings.GET("", settingsHandler.Get)
	settings.PATCH("", settingsHandler.Patch)

	sshHandler := handler.NewSshHandler(sshService)
	ssh := api.Group("/security/ssh")
	ssh.GET("", sshHandler.Get)
	ssh.POST("", sshHandler.Add)
	ssh.DELETE("/:fingerprint", sshHandler.Delete)
}

func startRouter() {
	url := config.Current.VertexURL()
	fmt.Printf("\n-- Vertex Client :: %s\n\n", url)
	log.Info("Vertex started", vlog.String("url", url))

	err := r.Start(fmt.Sprintf(":%s", config.Current.Port))
	if errors.Is(err, http.ErrServerClosed) {
		log.Info("Vertex closed")
	} else if err != nil {
		log.Error(err)
	}
}

func stopRouter() {
	// TODO: Stop() must also stop handleSignals()

	log.Info("gracefully stopping Vertex")

	ctx.DispatchEvent(types.EventServerStop{})

	cancelCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.Stop(cancelCtx)
	if err != nil {
		log.Error(err)
		return
	}
}
