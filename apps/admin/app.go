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

// docapi:admin title Vertex Admin
// docapi:admin description An admin panel to manage the Vertex platform.
// docapi:admin version 0.0.0
// docapi:admin filename admin

// docapi:admin url http://{ip}:{port-kernel}/api
// docapi:admin urlvar ip localhost The IP address of the server.
// docapi:admin urlvar port-kernel 7500 The port of the server.

// docapi:admin_kernel title Vertex Admin Kernel
// docapi:admin_kernel description An admin panel to manage the Vertex platform.
// docapi:admin_kernel version 0.0.0
// docapi:admin_kernel filename admin_kernel

// docapi:admin_kernel url http://{ip}:{port-kernel}/api
// docapi:admin_kernel urlvar ip localhost The IP address of the server.
// docapi:admin_kernel urlvar port-kernel 7501 The port of the server.

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

	settingsAdapter := adapter.NewAdminSettingsDbAdapter(db)
	hardwareAdapter := adapter.NewHardwareApiAdapter()
	sshAdapter := adapter.NewSshKernelApiAdapter()

	checksService := service.NewChecksService()
	settingsService := service.NewAdminSettingsService(settingsAdapter)
	hardwareService := service.NewHardwareService(hardwareAdapter)
	sshService := service.NewSshService(sshAdapter)

	service.NewNotificationsService(a.ctx, settingsAdapter)

	hardwareHandler := handler.NewHardwareHandler(hardwareService)
	hardware := r.Group("/hardware", middleware.Authenticated)
	// docapi:admin route /hardware/host get_host
	hardware.GET("/host", hardwareHandler.GetHost)
	// docapi:admin route /hardware/cpus get_cpus
	hardware.GET("/cpus", hardwareHandler.GetCPUs)
	// docapi:admin route /hardware/reboot reboot
	hardware.POST("/reboot", hardwareHandler.Reboot)

	sshHandler := handler.NewSshHandler(sshService)
	ssh := r.Group("/ssh", middleware.Authenticated)
	// docapi:admin route /ssh get_ssh_keys
	ssh.GET("", sshHandler.Get)
	// docapi:admin route /ssh add_ssh_key
	ssh.POST("", sshHandler.Add)
	// docapi:admin route /ssh delete_ssh_key
	ssh.DELETE("", sshHandler.Delete)
	// docapi:admin route /ssh/users get_ssh_users
	ssh.GET("/users", sshHandler.GetUsers)

	settingsHandler := handler.NewSettingsHandler(settingsService)
	settings := r.Group("/settings", middleware.Authenticated)
	// docapi:admin route /settings get_settings
	settings.GET("", settingsHandler.Get)
	// docapi:admin route /settings patch_settings
	settings.PATCH("", settingsHandler.Patch)

	updateHandler := handler.NewUpdateHandler(updateService, settingsService)
	update := r.Group("/update", middleware.Authenticated)
	// docapi:admin route /update get_updates
	update.GET("", updateHandler.Get)
	// docapi:admin route /update install_update
	update.POST("", updateHandler.Install)

	checksHandler := handler.NewChecksHandler(checksService)
	checks := r.Group("/admin/checks", middleware.Authenticated)
	// docapi:admin route /admin/checks admin_checks
	checks.GET("", apptypes.HeadersSSE, checksHandler.Check)

	return nil
}

func (a *App) InitializeKernel(r *router.Group) error {
	hardwareAdapter := adapter.NewHardwareKernelAdapter()
	sshAdapter := adapter.NewSshFsAdapter()

	hardwareService := service.NewHardwareKernelService(hardwareAdapter)
	sshService := service.NewSshKernelService(sshAdapter)

	hardwareHandler := handler.NewHardwareKernelHandler(hardwareService)
	hardware := r.Group("/hardware")
	// docapi:admin_kernel route /hardware/reboot reboot_kernel
	hardware.POST("/reboot", hardwareHandler.Reboot)

	sshHandler := handler.NewSshKernelHandler(sshService)
	ssh := r.Group("/ssh")
	// docapi:admin_kernel route /ssh get_ssh_keys_kernel
	ssh.GET("", sshHandler.Get)
	// docapi:admin_kernel route /ssh add_ssh_key_kernel
	ssh.POST("", sshHandler.Add)
	// docapi:admin_kernel route /ssh delete_ssh_key_kernel
	ssh.DELETE("", sshHandler.Delete)
	// docapi:admin_kernel route /ssh/users get_ssh_users_kernel
	ssh.GET("/users", sshHandler.GetUsers)

	return nil
}
