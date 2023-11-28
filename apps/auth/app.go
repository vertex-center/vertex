package auth

import (
	"github.com/vertex-center/vertex/apps/auth/adapter"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/service"
	"github.com/vertex-center/vertex/apps/auth/handler"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

var (
	authAdapter port.AuthAdapter

	AuthService port.AuthService
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
		ID:          "vx-auth",
		Name:        "Vertex Auth",
		Description: "Authentication app for Vertex",
		Icon:        "admin_panel_settings",
	}
}

func (a *App) Initialize(r *router.Group) error {
	authAdapter = adapter.NewAuthDbAdapter(a.ctx.Db())

	AuthService = service.NewAuthService(authAdapter)
	service.NewMigrationService(a.ctx)

	middleware.AuthService = AuthService

	authHandler := handler.NewAuthHandler(AuthService)
	auth := r.Group("/auth")
	// docapi:v route /app/vx-auth/auth/login auth_login
	auth.POST("/login", authHandler.Login)
	// docapi:v route /app/vx-auth/auth/register auth_register
	auth.POST("/register", authHandler.Register)
	// docapi:v route /app/vx-auth/auth/logout auth_logout
	auth.POST("/logout", middleware.Authenticated, authHandler.Logout)

	return nil
}
