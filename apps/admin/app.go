package admin

import (
	"path"

	"github.com/vertex-center/vertex/apps/admin/adapter"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/service"
	types2 "github.com/vertex-center/vertex/apps/admin/core/types"
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
	db, err := storage.NewDB(nil)
	if err != nil {
		return err
	}

	settingsAdapter := adapter.NewAdminSettingsDbAdapter(db)
	hardwareAdapter := adapter.NewHardwareApiAdapter()
	sshAdapter := adapter.NewSshKernelApiAdapter()

	settingsService := service.NewAdminSettingsService(settingsAdapter)
	dbService := service.NewDbService(a.ctx, db)
	hardwareService := service.NewHardwareService(hardwareAdapter)
	sshService := service.NewSshService(sshAdapter)

	service.NewNotificationsService(a.ctx, settingsAdapter)

	hardwareHandler := handler.NewHardwareHandler(hardwareService)
	hardware := r.Group("/hardware", middleware.Authenticated)
	// docapi:v route /app/admin/hardware/host get_host
	hardware.GET("/host", hardwareHandler.GetHost)
	// docapi:v route /app/admin/hardware/cpus get_cpus
	hardware.GET("/cpus", hardwareHandler.GetCPUs)
	// docapi:v route /app/admin/hardware/reboot reboot
	hardware.POST("/reboot", hardwareHandler.Reboot)

	sshHandler := handler.NewSshHandler(sshService)
	ssh := r.Group("/ssh", middleware.Authenticated)
	// docapi:v route /app/admin/ssh get_ssh_keys
	ssh.GET("", sshHandler.Get)
	// docapi:v route /app/admin/ssh add_ssh_key
	ssh.POST("", sshHandler.Add)
	// docapi:v route /app/admin/ssh delete_ssh_key
	ssh.DELETE("", sshHandler.Delete)
	// docapi:v route /app/admin/ssh/users get_ssh_users
	ssh.GET("/users", sshHandler.GetUsers)

	dbHandler := handler.NewDatabaseHandler(dbService)
	dbr := r.Group("/db", middleware.Authenticated)
	// docapi:v route /app/admin/db/dbms get_current_dbms
	dbr.GET("/dbms", dbHandler.GetCurrentDbms)
	// docapi:v route /app/admin/db/dbms migrate_to_dbms
	dbr.POST("/dbms", dbHandler.MigrateTo)

	settingsHandler := handler.NewSettingsHandler(settingsService)
	settings := r.Group("/settings", middleware.Authenticated)
	// docapi:v route /app/admin/settings get_settings
	settings.GET("", settingsHandler.Get)
	// docapi:v route /app/admin/settings patch_settings
	settings.PATCH("", settingsHandler.Patch)

	updateHandler := handler.NewUpdateHandler(updateService, settingsService)
	update := r.Group("/update", middleware.Authenticated)
	// docapi:v route /app/admin/update get_updates
	update.GET("", updateHandler.Get)
	// docapi:v route /app/admin/update install_update
	update.POST("", updateHandler.Install)

	return nil
}

func (a *App) InitializeKernel(r *router.Group) error {
	hardwareAdapter := adapter.NewHardwareKernelAdapter()
	sshAdapter := adapter.NewSshFsAdapter()

	hardwareService := service.NewHardwareKernelService(hardwareAdapter)
	sshService := service.NewSshKernelService(sshAdapter)

	hardwareHandler := handler.NewHardwareKernelHandler(hardwareService)
	hardware := r.Group("/hardware")
	// docapi:k route /app/admin/hardware/reboot reboot_kernel
	hardware.POST("/reboot", hardwareHandler.Reboot)

	sshHandler := handler.NewSshKernelHandler(sshService)
	ssh := r.Group("/ssh")
	// docapi:k route /app/admin/ssh get_ssh_keys_kernel
	ssh.GET("", sshHandler.Get)
	// docapi:k route /app/admin/ssh add_ssh_key_kernel
	ssh.POST("", sshHandler.Add)
	// docapi:k route /app/admin/ssh delete_ssh_key_kernel
	ssh.DELETE("", sshHandler.Delete)
	// docapi:k route /app/admin/ssh/users get_ssh_users_kernel
	ssh.GET("/users", sshHandler.GetUsers)

	return nil
}
