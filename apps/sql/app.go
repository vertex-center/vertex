package sql

import (
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/service"
	"github.com/vertex-center/vertex/apps/sql/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

var (
	sqlService port.SqlService
)

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
	return apptypes.Meta{
		ID:          "vx-sql",
		Name:        "Vertex SQL",
		Description: "Create and manage SQL databases.",
		Icon:        "database",
	}
}

func (a *App) Initialize(r *router.Group) error {
	sqlService = service.New(a.ctx)

	dbmsHandler := handler.NewDBMSHandler(sqlService)
	// docapi:v route /app/vx-sql/container/{container_uuid} vx_sql_get_dbms
	r.GET("/container/:container_uuid", middleware.Authenticated, dbmsHandler.Get)
	// docapi:v route /app/vx-sql/dbms/{dbms}/install vx_sql_install_dbms
	r.POST("/dbms/:dbms/install", middleware.Authenticated, dbmsHandler.Install)

	return nil
}
