package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	logsmeta "github.com/vertex-center/vertex/apps/logs/meta"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/server"
	"github.com/vertex-center/vertex/config"
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
	ctx *common.VertexContext
)

func main() {
	defer log.Default.Close()

	ensureNotRoot()

	about := common.NewAbout(version, commit, date)
	for _, a := range apps.Apps {
		meta := a.Meta()
		config.RegisterHost(meta.ID, meta.DefaultPort)
		config.RegisterHost(meta.ID+"-kernel", meta.DefaultKernelPort)
	}
	config.RegisterHost("vertex", "6130")
	config.ParseArgs(about)

	log.SetupAgent(*config.Current.Addr(logsmeta.Meta.ID))

	ctx = common.NewVertexContext(about, false)
	url := config.Current.Addr("vertex")

	info := openapi.Info{
		Title:       "Vertex",
		Description: "Create your self-hosted lab in one click.",
		Version:     ctx.About().Version,
	}

	go app.RunApps(about, apps.Apps)

	srv = server.New("main", &info, url, ctx)
	initRoutes(about)
	exitChan := srv.StartAsync()

	for err := range exitChan {
		if err != nil {
			log.Error(err)
		}
	}
}

func ensureNotRoot() {
	if os.Getuid() == 0 {
		log.Warn("while vertex-kernel must be run as root, the vertex user should not be root")
	}
}

func initRoutes(about common.About) {
	srv.Router.Engine().NoRoute(router.Handler(func(ctx *gin.Context) error {
		return errors.NewNotFound(nil, "route not found")
	}))

	a := srv.Router.Group("/api", "API", "Main API group", middleware.ReadAuth)
	a.GET("/about", []fizz.OperationOption{
		fizz.ID("getAbout"),
		fizz.Summary("Get server info"),
	}, router.Handler(func(ctx *gin.Context) (*common.About, error) {
		return &about, nil
	}))
}
