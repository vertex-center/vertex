package client

import (
	"context"
	"os"
	"path"

	"github.com/vertex-center/vertex/apps/client/meta"
	"github.com/vertex-center/vertex/apps/client/updates"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/common/updater"
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
		err := updater.Execute(context.Background(), ctx.About().Channel(),
			updates.NewVertexClientUpdater(path.Join(storage.FSPath, "client")),
		)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

func (a *App) Meta() appmeta.Meta {
	return meta.Meta
}
