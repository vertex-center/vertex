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
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/updates"
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

func (a *App) Initialize(r *router.Group) error {
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

	hardware.GET("/host", hardwareHandler.GetHostInfo(), hardwareHandler.GetHost)
	hardware.GET("/cpus", hardwareHandler.GetCPUsInfo(), hardwareHandler.GetCPUs)
	hardware.POST("/reboot", hardwareHandler.RebootInfo(), hardwareHandler.Reboot)

	ssh.GET("", sshHandler.GetInfo(), sshHandler.Get)
	ssh.POST("", sshHandler.AddInfo(), sshHandler.Add)
	ssh.DELETE("", sshHandler.DeleteInfo(), sshHandler.Delete)
	ssh.GET("/users", sshHandler.GetUsersInfo(), sshHandler.GetUsers)

	settings.GET("", settingsHandler.GetInfo(), settingsHandler.Get)
	settings.PATCH("", settingsHandler.PatchInfo(), settingsHandler.Patch)

	update.GET("", updateHandler.GetInfo(), updateHandler.Get)
	update.POST("", updateHandler.InstallInfo(), updateHandler.Install)

	checks.GET("", checksHandler.CheckInfo(), apptypes.HeadersSSE, checksHandler.Check)

	return nil
}

func (a *App) InitializeKernel(r *router.Group) error {
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

	hardware.POST("/reboot", hardwareHandler.RebootInfo(), hardwareHandler.Reboot)

	ssh.GET("", sshHandler.GetInfo(), sshHandler.Get)
	ssh.POST("", sshHandler.AddInfo(), sshHandler.Add)
	ssh.DELETE("", sshHandler.DeleteInfo(), sshHandler.Delete)
	ssh.GET("/users", sshHandler.GetUsersInfo(), sshHandler.GetUsers)

	return nil
}
