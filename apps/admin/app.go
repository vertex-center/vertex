package admin

import (
	"path"

	"github.com/vertex-center/vertex/apps/admin/adapter"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/service"
	types2 "github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/apps/admin/database"
	"github.com/vertex-center/vertex/apps/admin/handler"
	"github.com/vertex-center/vertex/apps/admin/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/core/types/storage"
	"github.com/vertex-center/vertex/updates"
	"github.com/wI2L/fizz"
)

var updateService port.UpdateService

type App struct {
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx

	if !ctx.Kernel() {
		baselinesApiAdapter := adapter.NewBaselinesApiAdapter()
		updateService = service.NewUpdateService(a.ctx, baselinesApiAdapter, []types2.Updater{
			updates.NewVertexUpdater(a.ctx.About()),
			updates.NewVertexClientUpdater(path.Join(storage.FSPath, "client")),
			updates.NewRepositoryUpdater("vertex_services", path.Join(storage.FSPath, "services"), "vertex-center", "services"),
		})
	}
}

func (a *App) Meta() apptypes.Meta {
	return meta.Meta
}

func (a *App) Initialize(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

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
		hardwareAdapter = adapter.NewHardwareApiAdapter()
		sshAdapter      = adapter.NewSshKernelApiAdapter()

		checksService   = service.NewChecksService()
		settingsService = service.NewSettingsService(settingsAdapter)
		hardwareService = service.NewHardwareService(hardwareAdapter)
		sshService      = service.NewSshService(sshAdapter)
		_               = service.NewNotificationsService(a.ctx, settingsAdapter)

		hardwareHandler = handler.NewHardwareHandler(hardwareService)
		sshHandler      = handler.NewSshHandler(sshService)
		settingsHandler = handler.NewSettingsHandler(settingsService)
		updateHandler   = handler.NewUpdateHandler(updateService, settingsService)
		checksHandler   = handler.NewChecksHandler(checksService)

		hardware = r.Group("/hardware", "Hardware", "", middleware.Authenticated)
		ssh      = r.Group("/ssh", "SSH", "", middleware.Authenticated)
		settings = r.Group("/settings", "Settings", "", middleware.Authenticated)
		update   = r.Group("/update", "Update", "", middleware.Authenticated)
		checks   = r.Group("/admin/checks", "Checks", "", middleware.Authenticated)
	)

	hardware.GET("/host", []fizz.OperationOption{
		fizz.ID("getHost"),
		fizz.Summary("Get host"),
		fizz.Description("Get host information."),
	}, hardwareHandler.GetHost())

	hardware.GET("/cpus", []fizz.OperationOption{
		fizz.ID("getCPUs"),
		fizz.Summary("Get CPUs"),
		fizz.Description("Get CPUs information."),
	}, hardwareHandler.GetCPUs())

	hardware.POST("/reboot", []fizz.OperationOption{
		fizz.ID("reboot"),
		fizz.Summary("Reboot"),
		fizz.Description("Reboot the host."),
	}, hardwareHandler.Reboot())

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
	}, apptypes.HeadersSSE(), checksHandler.Check())

	return nil
}

func (a *App) InitializeKernel(r *fizz.RouterGroup) error {
	var (
		hardwareAdapter = adapter.NewHardwareKernelAdapter()
		sshAdapter      = adapter.NewSshFsAdapter()

		hardwareService = service.NewHardwareKernelService(hardwareAdapter)
		sshService      = service.NewSshKernelService(sshAdapter)

		hardwareHandler = handler.NewHardwareKernelHandler(hardwareService)
		sshHandler      = handler.NewSshKernelHandler(sshService)

		hardware = r.Group("/hardware", "Hardware", "")
		ssh      = r.Group("/ssh", "SSH", "")
	)

	hardware.POST("/reboot", []fizz.OperationOption{
		fizz.ID("reboot"),
		fizz.Summary("Reboot"),
		fizz.Description("Reboot the host."),
	}, hardwareHandler.Reboot())

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
