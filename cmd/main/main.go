package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/adapter"
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/serviceeditor"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/handler"
	"github.com/vertex-center/vertex/migration"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/updates"
	"github.com/vertex-center/vlog"
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

	appsService          port.AppsService
	debugService         port.DebugService
	notificationsService service.NotificationsService
	hardwareService      port.HardwareService
	settingsService      port.SettingsService
	sshService           port.SshService
	updateService        port.UpdateService
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

	about := types.About{
		Version: version,
		Commit:  commit,
		Date:    date,

		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	initRouter()
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
	var (
		flagVersion        = flag.Bool("version", false, "Print vertex version")
		flagV              = flag.Bool("v", false, "Print vertex version")
		flagDate           = flag.Bool("date", false, "Print the release date")
		flagCommit         = flag.Bool("commit", false, "Print the commit hash")
		flagHost           = flag.String("host", config.Current.Host, "The Vertex access url")
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

func initRouter() {
	gin.SetMode(gin.ReleaseMode)
	ctx = types.NewVertexContext()
	r = router.New()
	r.Use(cors.Default())
	r.Use(ginutils.ErrorHandler())
	r.Use(ginutils.Logger("MAIN"))
	r.Use(gin.Recovery())
}

func initAdapters() {
	settingsFSAdapter = adapter.NewSettingsFSAdapter(nil)
	sshKernelApiAdapter = adapter.NewSshKernelApiAdapter()
	baselinesApiAdapter = adapter.NewBaselinesApiAdapter()
}

func initServices(about types.About) {
	// Update service must be initialized before all other services, because it
	// is responsible for downloading dependencies for other services.
	updateService = service.NewUpdateService(ctx, baselinesApiAdapter, []types.Updater{
		updates.NewVertexUpdater(about),
		updates.NewVertexClientUpdater(path.Join(storage.Path, "client")),
		updates.NewRepositoryUpdater("vertex_services", path.Join(storage.Path, "services"), "vertex-center", "services"),
	})
	appsService = service.NewAppsService(ctx, false, r,
		[]app.Interface{
			sql.NewApp(),
			tunnels.NewApp(),
			monitoring.NewApp(),
			containers.NewApp(),
			reverseproxy.NewApp(),
			serviceeditor.NewApp(),
		},
	)
	debugService = service.NewDebugService(ctx)
	notificationsService = service.NewNotificationsService(ctx, settingsFSAdapter)
	settingsService = service.NewSettingsService(settingsFSAdapter)
	//services.NewSetupService(r.ctx)
	hardwareService = service.NewHardwareService()
	sshService = service.NewSshService(sshKernelApiAdapter)
}

func initRoutes(about types.About) {
	// docapi:v title Vertex
	// docapi:v description A platform to manage your self-hosted server.
	// docapi:v version 0.0.0
	// docapi:v filename vertex

	// docapi:v url http://{ip}:{port}/api
	// docapi:v urlvar ip localhost The IP address of the server.
	// docapi:v urlvar port 6130 The port of the server.

	// docapi code 200 Success
	// docapi code 201 Created
	// docapi code 204 No content
	// docapi code 400 {Error} Bad request
	// docapi code 404 {Error} Resource not found
	// docapi code 409 {Error} Conflict
	// docapi code 422 {Error} Unprocessable entity
	// docapi code 500 {Error} Internal error

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
		debugHandler := handler.NewDebugHandler(debugService)
		debug := api.Group("/debug")
		// docapi:v route /debug/hard-reset hard_reset
		debug.POST("/hard-reset", debugHandler.HardReset)
	}

	appsHandler := handler.NewAppsHandler(appsService)
	apps := api.Group("/apps")
	// docapi:v route /apps get_apps
	apps.GET("", appsHandler.Get)

	hardwareHandler := handler.NewHardwareHandler(hardwareService)
	hardware := api.Group("/hardware")
	// docapi:v route /hardware get_hardware
	hardware.GET("", hardwareHandler.Get)

	updateHandler := handler.NewUpdateHandler(updateService, settingsService)
	update := api.Group("/update")
	// docapi:v route /update get_updates
	update.GET("", updateHandler.Get)
	// docapi:v route /update install_update
	update.POST("", updateHandler.Install)

	settingsHandler := handler.NewSettingsHandler(settingsService)
	settings := api.Group("/settings")
	// docapi:v route /settings get_settings
	settings.GET("", settingsHandler.Get)
	// docapi:v route /settings patch_settings
	settings.PATCH("", settingsHandler.Patch)

	sshHandler := handler.NewSshHandler(sshService)
	ssh := api.Group("/security/ssh")
	// docapi:v route /security/ssh get_ssh_keys
	ssh.GET("", sshHandler.Get)
	// docapi:v route /security/ssh add_ssh_key
	ssh.POST("", sshHandler.Add)
	// docapi:v route /security/ssh/{fingerprint} delete_ssh_key
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
