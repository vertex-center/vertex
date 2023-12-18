package client

import (
	"os"
	"path"

	"github.com/gin-contrib/static"
	"github.com/vertex-center/vertex/apps/client/meta"
	"github.com/vertex-center/vertex/apps/client/updates"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/common/updater"
	"github.com/wI2L/fizz"
)

type App struct {
	ctx *app.Context
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
	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.FSPath, "client", "dist"), true)))
	return nil
}
