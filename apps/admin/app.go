package admin

import (
	"github.com/vertex-center/vertex/apps/admin/core/service"
	"github.com/vertex-center/vertex/apps/admin/handler"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
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
		ID:          "admin",
		Name:        "Vertex Admin",
		Description: "Administer Vertex",
		Icon:        "admin_panel_settings",
	}
}

func (a *App) Initialize(r *router.Group) error {
	hardwareService := service.NewHardwareService()

	hardwareHandler := handler.NewHardwareHandler(hardwareService)
	hardware := r.Group("/hardware", middleware.Authenticated)
	// docapi:v route /app/admin/hardware/host get_host
	hardware.GET("/host", hardwareHandler.GetHost)
	// docapi:v route /app/admin/hardware/cpus get_cpus
	hardware.GET("/cpus", hardwareHandler.GetCPUs)

	return nil
}
