package sql

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/service"
	"github.com/vertex-center/vertex/apps/sql/handler"
	"github.com/vertex-center/vertex/common/app"
	"github.com/wI2L/fizz"
)

var (
	sqlService port.SqlService
)

var Meta = app.Meta{
	ID:          "sql",
	Name:        "Vertex SQL",
	Description: "Create and manage SQL databases.",
	Icon:        "database",
	DefaultPort: "7512",
	Dependencies: []*app.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
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

func (a *App) Meta() app.Meta {
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
