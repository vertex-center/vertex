package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/admin"
	"github.com/vertex-center/vertex/apps/auth"
	"github.com/vertex-center/vertex/apps/auth/middleware"
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
	"github.com/vertex-center/vertex/core/types/server"
	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/handler"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

// goreleaser will override version, commit and date
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	srv *server.Server
	ctx *types.VertexContext

	appsService  port.AppsService
	debugService port.DebugService
)

func main() {
	defer log.Default.Close()

	log.Info("Vertex starting...")

	parseArgs()

	checkNotRoot()

	about := types.About{
		Version: version,
		Commit:  commit,
		Date:    date,

		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	ctx = types.NewVertexContext(about, false)
	url := config.Current.URL("vertex")

	info := openapi.Info{
		Title:       "Vertex",
		Description: "Create your self-hosted lab in one click.",
		Version:     ctx.About().Version,
	}

	srv = server.New("main", &info, url, ctx)
	initServices()
	initRoutes(about)

	srv.Router.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.FSPath, "client", "dist"), true)))

	exitChan := srv.StartAsync()

	for err := range exitChan {
		if err != nil {
			log.Error(err)
		}
	}
}

func parseArgs() {
	var (
		flagVersion = flag.Bool("version", false, "Print vertex version")
		flagV       = flag.Bool("v", false, "Print vertex version")
		flagDate    = flag.Bool("date", false, "Print the release date")
		flagCommit  = flag.Bool("commit", false, "Print the commit hash")
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
}

func checkNotRoot() {
	if os.Getuid() == 0 {
		log.Warn("while vertex-kernel must be run as root, the vertex user should not be root")
	}
}

func initServices() {
	// Update service must be initialized before all other services, because it
	// is responsible for downloading dependencies for other services.
	appsService = service.NewAppsService(ctx, false, []app.Interface{
		admin.NewApp(),
		auth.NewApp(),
		sql.NewApp(),
		tunnels.NewApp(),
		monitoring.NewApp(),
		containers.NewApp(),
		reverseproxy.NewApp(),
		serviceeditor.NewApp(),
	})
	debugService = service.NewDebugService(ctx)
}

func initRoutes(about types.About) {
	srv.Router.Engine().NoRoute(router.Handler(func(c *gin.Context) error {
		return errors.NewNotFound(nil, "route not found")
	}))

	a := srv.Router.Group("/api", "API", "Main API group", middleware.ReadAuth())
	a.GET("/about", []fizz.OperationOption{
		fizz.ID("getAbout"),
		fizz.Summary("Get server info"),
	}, router.Handler(func(c *gin.Context) (*types.About, error) {
		return &about, nil
	}))

	if config.Current.Debug() {
		debugHandler := handler.NewDebugHandler(debugService)
		debug := a.Group("/debug", "Debug", "Routes only available with DEBUG=1", middleware.Authenticated())
		debug.POST("/hard-reset", debugHandler.HardResetInfo(), router.Handler(debugHandler.HardReset))
	}

	appsHandler := handler.NewAppsHandler(appsService)
	apps := a.Group("/apps", "Apps", "Apps", middleware.Authenticated())
	apps.GET("", appsHandler.GetAppsInfo(), router.Handler(appsHandler.GetApps))
}
