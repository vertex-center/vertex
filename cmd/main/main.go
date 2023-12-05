package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
	"github.com/vertex-center/vertex/handler"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/storage"
)

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
// docapi code 304 Not modified
// docapi code 400 {Error} Bad request
// docapi code 401 {Error} Unauthorized
// docapi code 404 {Error} Resource not found
// docapi code 409 {Error} Conflict
// docapi code 422 {Error} Unprocessable entity
// docapi code 500 {Error} Internal error

// goreleaser will override version, commit and date
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	srv *server.Server
	ctx *types.VertexContext

	appsService   port.AppsService
	debugService  port.DebugService
	checksService port.ChecksService
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
	addr := fmt.Sprintf(":%s", config.Current.Ports["VERTEX"])

	srv = server.New("main", addr, ctx)
	initServices()
	initRoutes(about)

	srv.Router.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.Path, "client", "dist"), true)))

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
		flagHost    = flag.String("host", config.Current.Host, "The Vertex access url")
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
	checksService = service.NewChecksService()
}

func initRoutes(about types.About) {
	srv.Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, router.Error{
			Code:          "resource_not_found",
			PublicMessage: "Resource not found.",
		})
	})

	a := srv.Router.Group("/api", middleware.ReadAuth)
	a.GET("/about", func(c *router.Context) {
		c.JSON(about)
	})

	if config.Current.Debug() {
		debugHandler := handler.NewDebugHandler(debugService)
		debug := a.Group("/debug", middleware.Authenticated)
		// docapi:v route /debug/hard-reset hard_reset
		debug.POST("/hard-reset", debugHandler.HardReset)
	}

	appsHandler := handler.NewAppsHandler(appsService)
	apps := a.Group("/apps", middleware.Authenticated)
	// docapi:v route /apps get_apps
	apps.GET("", appsHandler.Get)

	checksHandler := handler.NewChecksHandler(checksService)
	checks := a.Group("/admin/checks", middleware.Authenticated)
	// docapi:v route /admin/checks admin_checks
	checks.GET("", app.HeadersSSE, checksHandler.Check)
}
