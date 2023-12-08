package sql

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/sql/core/service"
	"github.com/vertex-center/vertex/apps/sql/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/wI2L/fizz"
)

var Meta = apptypes.Meta{
	ID:          "sql",
	Name:        "Vertex SQL",
	Description: "Create and manage SQL databases.",
	Icon:        "database",
	DefaultPort: "7512",
	Dependencies: []*apptypes.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
	},
}

type App struct {
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return Meta
}

func (a *App) Initialize(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth())

	var (
		sqlService  = service.New(a.ctx)
		dbmsHandler = handler.NewDBMSHandler(sqlService)
	)

	r.GET("/container/:container_uuid", []fizz.OperationOption{
		fizz.ID("getDBMS"),
		fizz.Summary("Get an installed DBMS"),
	}, middleware.Authenticated(), dbmsHandler.Get())

	r.POST("/dbms/:dbms/install", []fizz.OperationOption{
		fizz.ID("installDBMS"),
		fizz.Summary("Install a DBMS"),
	}, middleware.Authenticated(), dbmsHandler.Install())

	return nil
}
