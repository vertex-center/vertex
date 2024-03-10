package sql

import (
	authmeta "github.com/vertex-center/vertex/server/apps/auth/meta"
	"github.com/vertex-center/vertex/server/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/server/apps/containers/meta"
	logsmeta "github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/apps/sql/core/port"
	"github.com/vertex-center/vertex/server/apps/sql/core/service"
	"github.com/vertex-center/vertex/server/apps/sql/handler"
	"github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/common/app/appmeta"
	"github.com/wI2L/fizz"
)

var (
	sqlService port.SqlService
)

var Meta = appmeta.Meta{
	ID:          "sql",
	Name:        "Vertex SQL",
	Description: "Create and manage SQL databases.",
	Icon:        "database",
	DefaultPort: "7512",
	Dependencies: []*appmeta.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
		&logsmeta.Meta,
	},
}

type App struct {
	ctx *app.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *app.Context) {
	a.ctx = ctx
}

func (a *App) Meta() appmeta.Meta {
	return Meta
}

func (a *App) Initialize() error {
	sqlService = service.New(a.ctx)
	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	var (
		dbmsHandler = handler.NewDBMSHandler(sqlService)

		container = r.Group("/container/:container_uuid", "Container", "", middleware.Authenticated)
		dbms      = r.Group("/dbms/:dbms", "DBMS", "", middleware.Authenticated)
	)

	container.GET("", []fizz.OperationOption{
		fizz.ID("getDBMS"),
		fizz.Summary("Get an installed DBMS"),
	}, dbmsHandler.Get())

	dbms.POST("/install", []fizz.OperationOption{
		fizz.ID("installDBMS"),
		fizz.Summary("Install a DBMS"),
	}, dbmsHandler.Install())

	return nil
}
