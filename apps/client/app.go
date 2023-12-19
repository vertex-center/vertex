package client

import (
	"context"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/client/meta"
	"github.com/vertex-center/vertex/apps/client/updates"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/common/updater"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
	"github.com/wI2L/fizz"
)

type App struct {
	ctx    *app.Context
	router *router.Router
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *app.Context) {
	a.ctx = ctx

	if !ctx.Kernel() {
		bl, err := ctx.About().Baseline()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = updater.Install(bl, updates.NewVertexClientUpdater(path.Join(storage.FSPath, "client")))
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

func (a *App) Meta() appmeta.Meta {
	return meta.Meta
}

func (a *App) Initialize() error {
	a.router = router.New(nil,
		router.WithMiddleware(static.Serve("/", static.LocalFile(path.Join(".", storage.FSPath, "client", "dist"), true))),
	)

	a.router.GET("/ping", nil, func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	go func() {
		log.Info("client app starting", vlog.String("port", "6132"))
		err := a.router.Start(":6132")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.Info("client app stopped")
	}()

	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	return nil
}

func (a *App) Uninitialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	return a.router.Stop(ctx)
}
