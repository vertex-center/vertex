package admin

import (
	"github.com/vertex-center/vertex/apps/admin/adapter"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/service"
	"github.com/vertex-center/vertex/apps/admin/database"
	"github.com/vertex-center/vertex/apps/admin/handler"
	"github.com/vertex-center/vertex/apps/admin/meta"
	authmiddleware "github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/middleware"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/wI2L/fizz"
)

var (
	updateService   port.UpdateService
	checksService   port.ChecksService
	settingsService port.SettingsService
)

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
	return meta.Meta
}

func (a *App) Initialize() error {
	db, err := storage.NewDB(storage.DBParams{
		ID:         meta.Meta.ID,
		SchemaFunc: database.GetSchema,
		Migrations: database.Migrations,
	})
	if err != nil {
		return err
	}

	var (
		settingsAdapter = adapter.NewSettingsDbAdapter(db)
	)

	checksService = service.NewChecksService()
	settingsService = service.NewSettingsService(settingsAdapter)
	updateService = service.NewUpdateService(a.ctx)
	_ = service.NewNotificationsService(a.ctx, settingsAdapter)

	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(authmiddleware.ReadAuth)

	var (
		settingsHandler = handler.NewSettingsHandler(settingsService)
		updateHandler   = handler.NewUpdateHandler(updateService, settingsService)
		checksHandler   = handler.NewChecksHandler(checksService)

		settings = r.Group("/settings", "Settings", "", authmiddleware.Authenticated)
		update   = r.Group("/update", "Update", "", authmiddleware.Authenticated)
		checks   = r.Group("/admin/checks", "Checks", "", authmiddleware.Authenticated)
	)

	settings.GET("", []fizz.OperationOption{
		fizz.ID("getSettings"),
		fizz.Summary("Get settings"),
	}, settingsHandler.Get())

	settings.PATCH("", []fizz.OperationOption{
		fizz.ID("patchSettings"),
		fizz.Summary("Patch settings"),
	}, settingsHandler.Patch())

	update.GET("", []fizz.OperationOption{
		fizz.ID("getUpdate"),
		fizz.Summary("Get the latest update information"),
	}, updateHandler.Get())

	checks.GET("", []fizz.OperationOption{
		fizz.ID("check"),
		fizz.Summary("Get all checks"),
		fizz.Description("Check that all vertex requirements are met."),
	}, middleware.SSE, checksHandler.Check())

	return nil
}
