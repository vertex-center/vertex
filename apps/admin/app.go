package admin

import (
	"github.com/vertex-center/vertex/apps/admin/adapter"
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
	sshAdapter := adapter.NewSshKernelApiAdapter()

	hardwareService := service.NewHardwareService()
	sshService := service.NewSshService(sshAdapter)

	hardwareHandler := handler.NewHardwareHandler(hardwareService)
	hardware := r.Group("/hardware", middleware.Authenticated)
	// docapi:v route /app/admin/hardware/host get_host
	hardware.GET("/host", hardwareHandler.GetHost)
	// docapi:v route /app/admin/hardware/cpus get_cpus
	hardware.GET("/cpus", hardwareHandler.GetCPUs)

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

	return nil
}

func (a *App) InitializeKernel(r *router.Group) error {
	sshAdapter := adapter.NewSshFsAdapter()

	sshService := service.NewSshKernelService(sshAdapter)

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
