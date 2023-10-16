package sql

import (
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/service"
	"github.com/vertex-center/vertex/apps/sql/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

const (
	AppRoute = "/vx-sql"
)

var (
	sqlService port.SqlService
)

type App struct {
	*apptypes.App
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *apptypes.App) error {
	a.App = app

	c := app.Context()

	sqlService = service.New(c)

	app.Register(apptypes.Meta{
		ID:          "vx-sql",
		Name:        "Vertex SQL",
		Description: "Create and manage SQL databases.",
		Icon:        "database",
	})

	app.RegisterRoutes(AppRoute, func(r *router.Group) {
		dbmsHandler := handler.NewDBMSHandler(sqlService)
		r.GET("/container/:container_uuid", dbmsHandler.Get)
		r.POST("/dbms/:dbms/install", dbmsHandler.Install)
	})

	return nil
}
