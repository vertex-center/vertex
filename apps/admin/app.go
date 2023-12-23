package admin

import (
	"github.com/vertex-center/vertex/apps/admin/adapter"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/service"
	"github.com/vertex-center/vertex/apps/admin/database"
	"github.com/vertex-center/vertex/apps/admin/handler"
	"github.com/vertex-center/vertex/apps/admin/meta"
	"github.com/vertex-center/vertex/apps/admin/updates"
	authmiddleware "github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/middleware"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/common/updater"
	"github.com/wI2L/fizz"
)

var (
	updateService   port.UpdateService
	checksService   port.ChecksService
	settingsService port.SettingsService
	sshService      port.SshService

	sshKernelService port.SshKernelService
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
		updateService = service.NewUpdateService(a.ctx, []updater.Updater{
			updates.NewVertexUpdater(a.ctx.About()),
		})
	}
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
		sshAdapter      = adapter.NewSshKernelApiAdapter()
	)

	checksService = service.NewChecksService()
	settingsService = service.NewSettingsService(settingsAdapter)
	sshService = service.NewSshService(sshAdapter)
	_ = service.NewNotificationsService(a.ctx, settingsAdapter)

	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(authmiddleware.ReadAuth)

	var (
		sshHandler      = handler.NewSshHandler(sshService)
		settingsHandler = handler.NewSettingsHandler(settingsService)
		updateHandler   = handler.NewUpdateHandler(updateService, settingsService)
		checksHandler   = handler.NewChecksHandler(checksService)

		ssh      = r.Group("/ssh", "SSH", "", authmiddleware.Authenticated)
		settings = r.Group("/settings", "Settings", "", authmiddleware.Authenticated)
		update   = r.Group("/update", "Update", "", authmiddleware.Authenticated)
		checks   = r.Group("/admin/checks", "Checks", "", authmiddleware.Authenticated)
	)

	ssh.GET("", []fizz.OperationOption{
		fizz.ID("getSSHKeys"),
		fizz.Summary("Get all SSH keys"),
	}, sshHandler.Get())

	ssh.POST("", []fizz.OperationOption{
		fizz.ID("addSSHKey"),
		fizz.Summary("Add an SSH key"),
	}, sshHandler.Add())

	ssh.DELETE("", []fizz.OperationOption{
		fizz.ID("deleteSSHKey"),
		fizz.Summary("Delete SSH key"),
	}, sshHandler.Delete())

	ssh.GET("/users", []fizz.OperationOption{
		fizz.ID("getSSHUsers"),
		fizz.Summary("Get all users that can have SSH keys"),
	}, sshHandler.GetUsers())

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

	update.POST("", []fizz.OperationOption{
		fizz.ID("installUpdate"),
		fizz.Summary("Install the latest version"),
		fizz.Description("This endpoint will install the latest version of Vertex."),
	}, updateHandler.Install())

	checks.GET("", []fizz.OperationOption{
		fizz.ID("check"),
		fizz.Summary("Get all checks"),
		fizz.Description("Check that all vertex requirements are met."),
	}, middleware.SSE, checksHandler.Check())

	return nil
}

func (a *App) InitializeKernel() error {
	sshAdapter := adapter.NewSshFsAdapter()
	sshKernelService = service.NewSshKernelService(sshAdapter)
	return nil
}

func (a *App) InitializeKernelRouter(r *fizz.RouterGroup) error {
	var (
		sshHandler = handler.NewSshKernelHandler(sshKernelService)
		ssh        = r.Group("/ssh", "SSH", "")
	)

	ssh.GET("", []fizz.OperationOption{
		fizz.ID("getSSHKeys"),
		fizz.Summary("Get all SSH keys"),
	}, sshHandler.Get())

	ssh.POST("", []fizz.OperationOption{
		fizz.ID("addSSHKey"),
		fizz.Summary("Add an SSH key to the authorized_keys file"),
	}, sshHandler.Add())

	ssh.DELETE("", []fizz.OperationOption{
		fizz.ID("deleteSSHKey"),
		fizz.Summary("Delete an SSH key from the authorized_keys file"),
	}, sshHandler.Delete())

	ssh.GET("/users", []fizz.OperationOption{
		fizz.ID("getSSHUsers"),
		fizz.Summary("Get all users that can have SSH keys"),
	}, sshHandler.GetUsers())

	return nil
}
